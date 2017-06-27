package controllers

import (
	"errors"
	"github.com/AlexanderSuv/go-blog-BE/db/authors"
	"github.com/gorilla/mux"
	"net/http"
)

type Authors []authors.Author

func (as *Authors) Get(w http.ResponseWriter, r *http.Request) {
	blogAuthors, err := authors.Get()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Can`t read authors file")
		return
	}

	respondWithJson(w, blogAuthors)
}

func (as *Authors) GetById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	author := &authors.Author{Id: id}

	if err := author.Get(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Can`t read authors file")
		return
	}

	respondWithJson(w, author)
}

func (as *Authors) Post(w http.ResponseWriter, r *http.Request) {
	var author authors.Author
	id := mux.Vars(r)["id"]

	if err := parseJson(r, &author); err != nil {
		respondWithError(w, http.StatusBadRequest, err, "Invalid requiest payload")
		return
	}

	if !isValidAuthorRequest(&author) {
		respondWithError(w, http.StatusBadRequest, errors.New("Invalid requiest payload"), "Invalid requiest payload")
		return
	}

	author.Id = id

	if err := author.Save(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Can`t save author")
		return
	}

	respondWithJson(w, author)
}

func isValidAuthorRequest(a *authors.Author) bool {
	isValid := false
	if a.Id == "" && a.Registered == 0 {
		isValid = true
	}
	return isValid
}
