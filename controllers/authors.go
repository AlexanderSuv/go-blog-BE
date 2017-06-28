package controllers

import (
	"net/http"

	"github.com/AlexanderSuv/go-blog-BE/db"
	"github.com/gorilla/mux"
)

type Authors []db.Author

func (as *Authors) Get(w http.ResponseWriter, r *http.Request) {
	blogAuthors, err := db.GetAuthors()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, blogAuthors)
}

func (as *Authors) GetById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	author := &db.Author{Id: id}

	if err := author.Get(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, author)
}

func (as *Authors) Put(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var author db.Author

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
	var author db.Author

	if err := parseJson(r, &author); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	if err := db.NewAuthor(&author); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, author)
}

func (as *Authors) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	author := db.Author{Id: id}

	if err := author.Delete(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, struct{}{})
}
