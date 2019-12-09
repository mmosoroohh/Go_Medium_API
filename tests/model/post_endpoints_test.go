package model

import (
	"github.com/mmosoroohh/Go_Medium_API/api/models"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)

func TestAllPosts(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error occurred while refreshing users & posts table %v\n", err)
	}
	_, _, err = seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Error occurred while seeding user and post table %v\n", err)
	}
	posts, err := post.AllPosts(server.DB)
	if err != nil {
		t.Errorf("Error occurred while fetching posts: %v\n", err)
		return
	}
	assert.Equal(t, len(*posts), 2)
}

func TestCreatePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error occurred while refreshing user & post table %v\n", err)
	}

	user, err := seedUser()
	if err != nil {
		log.Fatalf("Can't seed user, please try again %v\n", err)
	}

	newPost := models.Post{
		ID:       1,
		Title:    "This is a test title",
		Content:  "This is a test content",
		AuthorID: user.ID,
	}
	savePost, err := newPost.SavePost(server.DB)
	if err != nil {
		t.Errorf("Error occurred while fetching post: %v\n", err)
		return
	}
	assert.Equal(t, newPost.ID, savePost.ID)
	assert.Equal(t, newPost.Title, savePost.Title)
	assert.Equal(t, newPost.Content, savePost.Content)
	assert.Equal(t, newPost.AuthorID, savePost.AuthorID)
}

func TestSinglePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error occurred while refreshing user & post table: %v\n", err)
	}
	singlePost, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error occurred while seeding table")
	}
	foundPost, err := post.SinglePost(server.DB, singlePost.ID)
	if err != nil {
		t.Errorf("Error occurred while fetching a user: %v\n", err)
		return
	}
	assert.Equal(t, foundPost.ID, singlePost.ID)
	assert.Equal(t, foundPost.Title, singlePost.Title)
	assert.Equal(t, foundPost.Content, singlePost.Content)
}

func TestUpdatePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error occurred while refreshing user & post table: %v\n", err)
	}
	singlePost, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error occurred while seeding table")
	}
	updatePost := models.Post{
		ID:       1,
		Title:    "This is the updated test title",
		Content:  "This is the updated test content",
		AuthorID: singlePost.AuthorID,
	}
	updatedPost, err := updatePost.UpdatePost(server.DB)
	if err != nil {
		t.Errorf("Error occurred while updating a post: %v\n", err)
		return
	}
	assert.Equal(t, updatedPost.ID, updatePost.ID)
	assert.Equal(t, updatedPost.Title, updatePost.Title)
	assert.Equal(t, updatedPost.Content, updatePost.Content)
	assert.Equal(t, updatedPost.AuthorID, updatePost.AuthorID)
}

func TestDeletePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error occurred while refreshing user & post table: %v\n", err)
	}
	singlePost, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error occurred while seeding tables")
	}
	deletePost, err := post.DeletePost(server.DB, singlePost.ID, singlePost.AuthorID)
	if err != nil {
		t.Errorf("Error occurred while deleting the post: %v\n", err)
		return
	}
	assert.Equal(t, int(deletePost), 1)
	assert.Equal(t, deletePost, int64(1))
}
