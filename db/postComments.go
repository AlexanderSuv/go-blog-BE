package db

import (
	"errors"
	"github.com/nu7hatch/gouuid"
	"time"
)

// TYPES

type Comment struct {
	Id      string `json:"id" validate:"required"`
	Content string `json:"content" validate:"required"`
	Author  string `json:"author" validate:"required"`
	Updated int64  `json:"updated" validate:"required"`
}

// TYPE HELPERS

func softAssignComment(from *Comment, to *Comment) {
	if from.Content != "" {
		to.Content = from.Content
	}

	if from.Author != "" {
		to.Author = from.Author
	}
}

func getCommentById(id string, cs []*Comment) *Comment {
	var result *Comment
	for _, comment := range cs {
		if comment.Id == id {
			result = comment
		}
	}

	return result
}

func getCommentArrIndex(id string, comments []*Comment) int {
	result := -1
	for index, comment := range comments {
		if comment.Id == id {
			result = index
			break
		}
	}

	return result
}

// API

func (p *Post) GetComments(offset, limit int) ([]*Comment, error) {
	if err := p.Get(); err != nil {
		return nil, err
	}

	comments := make([]*Comment, len(p.Comments))
	copy(comments, p.Comments)

	return comments, nil
}

func (p *Post) GetCommentById(commentId string) (*Comment, error) {
	if err := p.Get(); err != nil {
		return nil, err
	}

	comment := getCommentById(commentId, p.Comments)
	if comment == nil {
		return nil, errors.New("No comment found")
	}
	return comment, nil
}

func (p *Post) NewComment(c *Comment) error {
	u, uuidErr := uuid.NewV4()
	if uuidErr != nil {
		return uuidErr
	}

	c.Id = u.String()
	c.Updated = time.Now().UnixNano() / 1000000

	if err := validate.Struct(c); err != nil {
		return err
	}

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

	posts[index].Comments = append(posts[index].Comments, c)
	return writePosts(posts)
}

func (p *Post) UpdateComment(c *Comment) error {
	postsFileMutex.Lock()
	defer postsFileMutex.Unlock()
	posts, err := readPosts()
	if err != nil {
		return err
	}
	postIndex := getPostArrIndex(p.Id, posts)
	*p = *posts[postIndex]

	commentIndex := getCommentArrIndex(c.Id, p.Comments)
	if commentIndex == -1 {
		return errors.New("No comment found")
	}

	commentToUpdate := p.Comments[commentIndex]
	softAssignComment(c, commentToUpdate)
	commentToUpdate.Updated = time.Now().UnixNano() / 1000000

	if err := validate.Struct(commentToUpdate); err != nil {
		return err
	}

	*c = *commentToUpdate
	return writePosts(posts)
}

func (p *Post) DeleteComment(id string) error {
	postsFileMutex.Lock()
	defer postsFileMutex.Unlock()
	posts, err := readPosts()
	if err != nil {
		return err
	}

	postIndex := getPostArrIndex(p.Id, posts)
	if postIndex == -1 {
		return errors.New("No post found")
	}

	commentIndex := getCommentArrIndex(id, posts[postIndex].Comments)
	if commentIndex == -1 {
		return errors.New("No comment found")
	}

	posts[postIndex].Comments = append(posts[postIndex].Comments[:commentIndex], posts[postIndex].Comments[commentIndex+1:]...)

	return writePosts(posts)
}
