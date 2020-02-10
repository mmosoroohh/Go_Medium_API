package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mmosoroohh/Go_Medium_API/api/models"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCreatePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	person, err := seedUser()
	if err != nil {
		log.Fatalf("Error Occurred seeding user %v\n", err)
	}
	token, err := server.SignIn(person.Email, "password") // Note password in the DB is already hashed, we need it unhashed.
	if err != nil {
		log.Fatalf("Error occurred login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		title        string
		content      string
		authorId     uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"title": "This is the title", "content": "This is the content", "author_id": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			title:        "This is the title",
			content:      "This is the content",
			authorId:     person.ID,
			errorMessage: "",
		},
		{
			// Passing Already Existing Title
			inputJSON:    `{"title": "This is the title", "content": "This is the content", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Exist",
		},
		{
			// No Token provided
			inputJSON:    `{"title": "This is title two", "content": "This is content two", "author_id":1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// Incorrect token provided
			inputJSON:    `{"title": "This is title three", "content": "This is content three", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "Incorrect token given",
			errorMessage: "Unauthorized",
		},
		{
			// Missing Title
			inputJSON:    `{"title": "", "content": "This content", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Title Required",
		},
		{
			// Missing Content
			inputJSON:    `{"title": "This title", "content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Content Required",
		},
		{
			// Missing Author
			inputJSON:    `{"title": "this title", "content": "this content"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Author Required",
		},
		{
			// User 2 uses 1 token
			inputJSON:    `{"title": "This is the awesome title", "content": "This is the awesome content", "author_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/posts", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("Error Occurred: %v\n", err)
		}

		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreatePost)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rec, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rec.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Error Occurred converting to json: %v", err)
		}
		assert.Equal(t, rec.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.authorId))
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetPosts(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Errorf("Error Occurred: %v\n", err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetPosts)
	handler.ServeHTTP(rec, req)

	var posts []models.Post
	err = json.Unmarshal([]byte(rec.Body.String()), &posts)

	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, len(posts), 2)
}

func TestSinglePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatal(err)
	}

	postSample := []struct {
		id           string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(post.ID)),
			statusCode: 200,
			title:      post.Title,
			content:    post.Content,
			author_id:  post.AuthorID,
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}

	for _, v := range postSample {

		req, err := http.NewRequest("GET", "/posts", nil)
		if err != nil {
			t.Errorf("Error Occurred : %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetPost)
		handler.ServeHTTP(rec, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rec.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Error Occurred converting to json: %v", err)
		}
		assert.Equal(t, rec.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, post.Title, responseMap["title"])
			assert.Equal(t, post.Content, responseMap["content"])
			assert.Equal(t, float64(post.AuthorID), responseMap["author_id"])
		}
	}
}

func TestUpdatePost(t *testing.T) {

	var UserEmail, UserPassword string
	var AuthorID uint32
	var PostID uint64

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	// Get the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		UserEmail = user.Email
		UserPassword = "password" //Note the password in the database is already hashed, we want unhashed.
	}

	//Login user and get authentication token
	token, err := server.SignIn(UserEmail, UserPassword)
	if err != nil {
		log.Fatalf("Error Occurred cannot login: %v", token)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get the first post
	for _, post := range posts {
		if post.ID == 2 {
			continue
		}
		PostID = post.ID
		AuthorID = post.AuthorID
	}

	// fmt.Printf("This is the auth post: %v\n", AuthPostID)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			//Convert int64 to int before converting to string
			id:           strconv.Itoa(int(PostID)),
			updateJSON:   `{"title": "Updated post", "content": "Updated content", "author_id":1}`,
			statusCode:   200,
			title:        "Updated post",
			content:      "Updated content",
			author_id:    AuthorID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			//No token provided
			id:           strconv.Itoa(int(PostID)),
			updateJSON:   `{"title": "Another Updated title", "content": "Another updated content", "author_id":1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			//Incorrect token provided
			id:           strconv.Itoa(int(PostID)),
			updateJSON:   `{"title": "Another title", "content": "Another content", "author_id": 1}`,
			tokenGiven:   "incorrect token given",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			//Note: "Title 2" belongs to post 2, and title must be unique
			id:           strconv.Itoa(int(PostID)),
			updateJSON:   `{"title":"Title 2", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			id:           strconv.Itoa(int(PostID)),
			updateJSON:   `{"title":"", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			id:           strconv.Itoa(int(PostID)),
			updateJSON:   `{"title":"Awesome title", "content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			id:           strconv.Itoa(int(PostID)),
			updateJSON:   `{"title":"This is another title", "content": "This is the updated content"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(PostID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/posts", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdatePost)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rec, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rec.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rec.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just to match the type of the json we receive thats why we used float64
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeletePost(t *testing.T) {

	var UserEmail, UserPassword string
	var UserID uint32
	var PostID uint64

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	//Let's get the Second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		UserEmail = user.Email
		UserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}

	//Login the user and get the authentication token
	token, err := server.SignIn(UserEmail, UserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get the second post
	for _, post := range posts {
		if post.ID == 1 {
			continue
		}
		PostID = post.ID
		UserID = post.AuthorID
	}
	postSample := []struct {
		id           string
		author_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(PostID)),
			author_id:    UserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// Empty token is passed
			id:           strconv.Itoa(int(PostID)),
			author_id:    UserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// Incorrect token is passed
			id:           strconv.Itoa(int(PostID)),
			author_id:    UserID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			author_id:    1,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range postSample {

		req, _ := http.NewRequest("GET", "/posts", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeletePost)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rec, req)

		assert.Equal(t, rec.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rec.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
