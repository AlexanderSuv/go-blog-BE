package controllers

import (
	"github.com/AlexanderSuv/go-blog-BE/db"
	"github.com/gorilla/mux"
	"net/http"
)

type Posts []db.Author

func (ps *Posts) Get(w http.ResponseWriter, r *http.Request) {
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

	posts, err := db.GetPosts(offset, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, posts)
}

func (ps *Posts) GetById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	post := &db.Post{Id: id}

	if err := post.Get(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, post)
}

func (ps *Posts) Post(w http.ResponseWriter, r *http.Request) {
	var post db.Post

	if err := parseJson(r, &post); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	if err := db.NewPost(&post); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, post)
}

func (ps *Posts) Put(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var post db.Post

	if err := parseJson(r, &post); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	post.Id = id
	if err := post.Update(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, post)
}

func (ps *Posts) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	post := db.Post{Id: id}

	if err := post.Delete(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, struct{}{})
}
