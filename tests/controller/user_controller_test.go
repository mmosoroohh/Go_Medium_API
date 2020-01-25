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

func TestCreateUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		inputJSON    string
		statusCode   int
		username     string
		email        string
		errorMessage string
	}{
		{
			inputJSON:    `{"username": "mmosoroohh", "email": "mmosoroohh@gmail.com", "password": "password"}`,
			statusCode:   201,
			username:     "mmosoroohh",
			email:        "mmosoroohh@gmail.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"username": "Arnold", "email": "mmosoroohh@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "Email Already Exists",
		},
		{
			inputJSON:    `{"username": "mmosoroohh", "email": "email@gmail", "password": "password"}`,
			statusCode:   500,
			errorMessage: "Username already Exists",
		},
		{
			inputJSON:    `{"username":"mmosoroohh", "email": "mmosoroohhgmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"username": "", "email": "mmosoroohh@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Username Required",
		},
		{
			inputJSON:    `{"username": "mmosoroohh", "email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Email Required",
		},
		{
			inputJSON:    `{"username": "mmosoroohh", "email": "mmosoroohh@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "Password Required",
		},
	}

	for _, v := range samples {

		request, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("Error Occurred: %v", err)
		}
		record := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateUser)
		handler.ServeHTTP(record, request)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(record.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Error occurred converting to json: %v", err)
		}
		assert.Equal(t, record.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["username"], v.username)
			assert.Equal(t, responseMap["email"], v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetUsers(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Errorf("Error occurred: %v\n", err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetUsers)
	handler.ServeHTTP(rec, req)

	var users []models.User
	err = json.Unmarshal([]byte(rec.Body.String()), &users)
	if err != nil {
		log.Fatal("Error occurred converting to json: %v\n", err)
	}
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, len(users), 2)
}

func TestGetUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	person, err := seedUser()
	if err != nil {
		log.Fatal(err)
	}

	sampleUser := []struct {
		id           string
		statusCode   int
		username     string
		email        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(person.ID)),
			statusCode: 200,
			username:   person.Username,
			email:      person.Email,
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}

	for _, v := range sampleUser {

		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("Error occurred: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetUser)
		handler.ServeHTTP(rec, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rec.Body.String()), &responseMap)
		if err != nil {
			log.Fatal("Error occurred converting to json: %v", err)
		}

		assert.Equal(t, rec.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, person.Username, responseMap["username"])
			assert.Equal(t, person.Email, responseMap["email"])
		}
	}
}

func TestUpdateUser(t *testing.T) {

	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	users, err := seedUsers() //we need atleast two users to properly check the update
	if err != nil {
		log.Fatal("Error occurred seeding the users: %v\n", err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
		AuthEmail = user.Email
		AuthPassword = "password"
		// Note the password is the database is already hashed, we want unhashed password
	}
	// Login the user and get the authentication token
	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatal("Error occurred cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id             string
		updateJSON     string
		statusCode     int
		updateUsername string
		updateEmail    string
		tokenGiven     string
		errorMessage   string
	}{
		{
			// Convert int32 to int first before converting to a string
			id:             strconv.Itoa(int(AuthID)),
			updateJSON:     `{"username": "mmosoroohh", "email": "mmosoroohh@gmail.com", "password": "password"}`,
			statusCode:     200,
			updateUsername: "mmosoroohh",
			updateEmail:    "mmosoroohh@gmail.com",
			tokenGiven:     tokenString,
			errorMessage:   "",
		},
		{
			// When password field is empty
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nickname": "joe", "email": "joe@gmail.com", "password": ""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Password Required",
		},
		{
			// When No token is provided
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username": "mary", "email": "mary@gmail.com", "password": "password"}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When an incorrect token is passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username": "brian", "email": "brian@gmail.com", "password": "password"}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			// Remember "joedoe@gmail.com" belongs to user two
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username": "joedoe", "email": "joedoe@gmail.com", "password": "password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Email Already Exist",
		},
		{
			// Remember "John Doe" belongs to user two
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username": "joedoe", "email": "mmosoroohh@gmail.com", "password": "password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Username Already Exists",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username": "mmosoroohh", "email": "mmosoroohh@gmail.com", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Invalid Email",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username": "", "email": "mmosoroohh@gmail.com", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Username Required",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username": "mmosoroohh", "email": "", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Email Required",
		},
		{
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("Error Occurred: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateUser)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rec, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rec.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Error Occurred converting to json: %v", err)
		}

		assert.Equal(t, rec.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["username"], v.updateUsername)
			assert.Equal(t, responseMap["email"], v.updateEmail)
		}

		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
