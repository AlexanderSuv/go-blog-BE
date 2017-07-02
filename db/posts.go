package db

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"encoding/json"
	"errors"
	"github.com/nu7hatch/gouuid"
	"io/ioutil"
	"time"
)

// INITIALIZATION

var pathToPosts = getPathToPosts()
var postsFileMutex = &sync.RWMutex{}

func getPathToPosts() string {
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
	Comments []*Comment `json:"comments" validate:"dive"`
}


// TYPE HELPERS

func readPosts() ([]*Post, error) {
	raw, readFileErr := ioutil.ReadFile(pathToPosts)
	if readFileErr != nil {
		fmt.Println(readFileErr)
		return nil, errors.New("Can`t read posts file")
	}

	var posts []*Post
	if parseErr := json.Unmarshal(raw, &posts); parseErr != nil {
		return nil, parseErr
	}

	return posts, nil
}

func writePosts(a []*Post) error {
	bytes, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pathToPosts, bytes, 0600)
}

func getPostArrIndex(id string, posts []*Post) int {
	result := -1
	for index, post := range posts {
		if post.Id == id {
			result = index
			break
		}
	}

	return result
}

func softAssignPost(from *Post, to *Post) {
	if from.Content != "" {
		to.Content = from.Content
	}

	if from.Author != "" {
		to.Author = from.Author
	}
}

// API

func NewPost(p *Post) error {
	u, uuidErr := uuid.NewV4()
	if uuidErr != nil {
		return uuidErr
	}

	p.Id = u.String()
	p.Updated = time.Now().UnixNano() / 1000000
	p.Comments = nil

	if err := validate.Struct(p); err != nil {
		return err
	}

	postsFileMutex.Lock()
	defer postsFileMutex.Unlock()

	posts, err := readPosts()
	if err != nil {
		return err
	}

	posts = append(posts, p)
	return writePosts(posts)
}

func GetPosts(offset int, limit int) ([]*Post, error) {
	postsFileMutex.RLock()
	defer postsFileMutex.RUnlock()

	posts, err := readPosts()
	if err != nil {
		return nil, err
	}

	result := make([]*Post, 0)
	length := len(posts)
	if offset > length {
		return result, nil
	}

	if offset+limit > length {
		limit = length - offset
	}

	for i := offset; i < offset+limit; i++ {
		result = append(result, posts[i])
	}

	return result, nil
}

func (p *Post) Get() error {
	postsFileMutex.Lock()
	defer postsFileMutex.Unlock()
	posts, err := readPosts()
	if err != nil {
		return err
	}
	index := getPostArrIndex(p.Id, posts)

	if index == -1 {
		return errors.New("No post found.")
	}

	*p = *posts[index]
	return nil
}

func (p *Post) Update() error {
	postsFileMutex.Lock()
	defer postsFileMutex.Unlock()

	posts, err := readPosts()
	if err != nil {
		return err
	}

	index := getPostArrIndex(p.Id, posts)

	if index == -1 {
		return errors.New("Post not found. Can`t save")
	}

	postToUpdate := posts[index]

	softAssignPost(p, postToUpdate)

	timeStamp := time.Now().UnixNano() / 1000000

	postToUpdate.Updated = timeStamp

	*p = *postToUpdate

	if err := validate.Struct(postToUpdate); err != nil {
		return err
	}
	return writePosts(posts)
}

func (p *Post) Delete() error {
	postsFileMutex.Lock()
	defer postsFileMutex.Unlock()

	posts, err := readPosts()
	if err != nil {
		return err
	}
	index := getPostArrIndex(p.Id, posts)

	if index == -1 {
		return errors.New("No post found. Can`t delete.")
	}

	posts = append(posts[:index], posts[index+1:]...)
	return writePosts(posts)
}
