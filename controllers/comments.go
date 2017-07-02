package controllers

import (
	"github.com/AlexanderSuv/go-blog-BE/db"
	"github.com/gorilla/mux"
	"net/http"
)

type Comments []*db.Comment

func (cs *Comments) Get(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	offset, err := queryStringToInt(qs, "offset")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	if offset < 0 {
		offset = defaultOffset
	}

	limit, err := queryStringToInt(qs, "limit")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	if limit < 0 {
		limit = defaultLimit
	}

	postId := mux.Vars(r)["postId"]
	post := &db.Post{Id: postId}

	if err := post.Get(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	comments, err := post.GetComments(offset, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, comments)
}

func (cs *Comments) GetById(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["postId"]
	post := &db.Post{Id: postId}

	commentId := mux.Vars(r)["commentId"]
	comment, err := post.GetCommentById(commentId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, comment)
}

func (cs *Comments) Post(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["postId"]
	post := &db.Post{Id: postId}

	var comment db.Comment
	if err := parseJson(r, &comment); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	if err := post.NewComment(&comment); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, comment)
}

func (cs *Comments) Put(w http.ResponseWriter, r *http.Request) {
	var comment db.Comment
	if err := parseJson(r, &comment); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	postId := mux.Vars(r)["postId"]
	post := &db.Post{Id: postId}

	commentId := mux.Vars(r)["commentId"]
	comment.Id = commentId

	if err := post.UpdateComment(&comment); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, comment)
}

func (cs *Comments) Delete(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["postId"]
	post := &db.Post{Id: postId}

	commentId := mux.Vars(r)["commentId"]

	if err := post.DeleteComment(commentId); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, struct{}{})
}
