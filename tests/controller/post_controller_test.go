package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mmosoroohh/Go_Medium_API/api/models"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"net/http"
	"net/http/httptest"
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
