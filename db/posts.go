package db

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"encoding/json"
	"errors"
	"io/ioutil"
)

// INITIALIZATION

var pathToPosts = getpathToPosts()
var postsFileMutex = &sync.RWMutex{}

func getpathToPosts() string {
	path, filenameErr := filepath.Abs("./db/db_data/posts.json")
	if filenameErr != nil {
		fmt.Println(filenameErr.Error())
		os.Exit(1)
	}
	return path
}

// TYPES

type Post struct {
	Id       string     `json:"id" validate:"required"`
	Content  string     `json:"content" validate:"required"`
	Updated  int64      `json:"updated" validate:"required"`
	Author   string     `json:"author" validate:"required"`
	Comments []*Comment `json:"comment" validate:"required,dive,required"`
}

type Comment struct {
	Content string `json:"content" validate:"required"`
	Author  string `json:"author" validate:"required"`
	Updated int64  `json:"updated" validate:"required"`
}

// TYPE HELPERS

func readPosts() (*[]Post, error) {
	raw, readFileErr := ioutil.ReadFile(pathToPosts)
	if readFileErr != nil {
		fmt.Println(readFileErr)
		return nil, errors.New("Can`t read posts file")
	}

	var posts []Post
	if parseErr := json.Unmarshal(raw, &posts); parseErr != nil {
		return nil, parseErr
	}

	return &posts, nil
}

func writePosts(a []*Post) error {
	bytes, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pathToPosts, bytes, 0600)
}

func getPostArrIndex(id string, posts *[]Post) int {
	result := -1
	for index, post := range *posts {
		if post.Id == id {
			result = index
			break
		}
	}

	return result
}

// API

func GetPosts() (*[]Post, error) {
	postsFileMutex.RLock()
	defer postsFileMutex.RUnlock()
	return readPosts()
}
