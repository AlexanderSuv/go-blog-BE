package controllers

import (
	"net/http"

	"github.com/AlexanderSuv/go-blog-BE/db/authors"
	"github.com/gorilla/mux"
)

type Authors []authors.Author

func (as *Authors) Get(w http.ResponseWriter, r *http.Request) {
	blogAuthors, err := authors.Get()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, blogAuthors)
}

func (as *Authors) GetById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	author := &authors.Author{Id: id}

	if err := author.Get(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, author)
}

func (as *Authors) Put(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var author authors.Author

	if err := parseJson(r, &author); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	author.Id = id
	if err := author.Update(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, author)
}

func (as *Authors) Post(w http.ResponseWriter, r *http.Request) {
	var author authors.Author

	if err := parseJson(r, &author); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	if err := authors.NewAuthor(&author); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, author)
}
