package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		panic(err)
	}
}

func respondWithError(w http.ResponseWriter, code int, err error, message string) {
	log.Println(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(map[string]string{"message": message}); err != nil {
		panic(err)
	}
}

func parseJson(r *http.Request, a interface{}) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	return decoder.Decode(a)
}
