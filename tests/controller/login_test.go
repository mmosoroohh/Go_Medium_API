package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignIn(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	person, err := seedUser()
	if err != nil {
		fmt.Printf("Error occurred %v\n", err)
	}

	sample := []struct {
		email        string
		password     string
		errorMessage string
	}{
		{
			email:        person.Email,
			password:     "password",
			errorMessage: "",
		},
		{
			email:        person.Email,
			password:     "Wrong password",
			errorMessage: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			email:        "Wrong email",
			password:     "password",
			errorMessage: "Record not found",
		},
	}

	for _, v := range sample {
		token, err := server.SignIn(v.email, v.password)
		if err != nil {
			assert.Equal(t, err, errors.New(v.errorMessage))
		} else {
			assert.NotEqual(t, token, "")
		}
	}
}

func TestLogin(t *testing.T) {

	refreshUserTable()

	_, err := seedUser()
	if err != nil {
		fmt.Printf("Error occurred %v\n", err)
	}

	samples := []struct {
		inputJSON    string
		statusCode   int
		email        string
		password     string
		errorMessage string
	}{
		{
			inputJSON:    `{"email": "joe@gmail.com", "password": "password"}`,
			statusCode:   200,
			errorMessage: "",
		},
		{
			inputJSON:    `{"email": "joe@gmail", "password": "wrong password"}`,
			statusCode:   422,
			errorMessage: "Incorrect Password",
		},
		{
			inputJSON:    `{"email": "doe@gmail", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Incorrect Details",
		},
		{
			inputJSON:    `{"email": "joegmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Email Required",
		},
		{
			inputJSON:    `{"email": "joe@gmail", "password": ""}`,
			statusCode:   422,
			errorMessage: "Password Required",
		},
	}

	for _, v := range samples {

		request, err := http.NewRequest("POST", "/login", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("Error occurred: %v", err)
		}
		record := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Login)
		handler.ServeHTTP(record, request)

		assert.Equal(t, record.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.NotEqual(t, record.Body.String(), "")
		}

		if v.statusCode == 422 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(record.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Error occurred converting to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
