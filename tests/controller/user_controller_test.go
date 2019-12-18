package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"net/http"
	"net/http/httptest"
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
