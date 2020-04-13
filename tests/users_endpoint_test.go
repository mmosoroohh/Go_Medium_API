package tests

import (
	"github.com/mmosoroohh/Go_Medium_API/api/models"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)

func TestAllUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	users, err := user.AllUsers(server.DB)
	if err != nil {
		t.Errorf("Error occurred with getting users: %v\n", err)
		return
	}
	assert.Equal(t, len(*users), 2)
}

func TestCreateUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	newUser := models.User{
		ID:       1,
		Username: "test",
		Email:    "test@gmail.com",
		Password: "password",
	}
	saveUser, err := newUser.SaveUser(server.DB)
	if err != nil {
		t.Errorf("Error occurred while saving saving the user: %v\n", err)
		return
	}
	assert.Equal(t, newUser.ID, saveUser.ID)
	assert.Equal(t, newUser.Username, saveUser.Username)
	assert.Equal(t, newUser.Email, saveUser.Email)
}

func TestSingleUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	person, err := seedUser()
	if err != nil {
		log.Fatalf("Can't seed user table: %v", err)
	}

	singleUser, err := user.SingleUser(server.DB, person.ID)
	if err != nil {
		t.Errorf("Error occurred while getting a user: %v\n", err)
		return
	}
	assert.Equal(t, singleUser.ID, person.ID)
	assert.Equal(t, singleUser.Email, person.Email)
	assert.Equal(t, singleUser.Username, person.Username)
}

func TestUpdateUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedUser()
	if err != nil {
		log.Fatalf("Can't seed user: %v\n", err)
	}

	updateUser := models.User{
		ID:       1,
		Username: "testUser",
		Email:    "testuser@gmail.com",
		Password: "password",
	}

	updatedUser, err := updateUser.UpdateAUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("Error occurred while updating user: %v\n", err)
		return
	}
	assert.Equal(t, updatedUser.ID, updateUser.ID)
	assert.Equal(t, updatedUser.Email, updateUser.Email)
	assert.Equal(t, updatedUser.Username, updateUser.Username)
}

func TestDeleteUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	person, err := seedUser()
	if err != nil {
		log.Fatalf("Can't seed user: %v\n", err)
	}

	deleteUser, err := user.DeleteUser(server.DB, person.ID)
	if err != nil {
		t.Errorf("Error occurred while deleting user: %v\n", err)
		return
	}
	assert.Equal(t, deleteUser, int64(1))
}
