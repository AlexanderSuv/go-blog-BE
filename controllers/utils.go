package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const defaultOffset = 0
const defaultLimit = 10

func respondWithJson(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		panic(err)
	}
}

func respondWithError(w http.ResponseWriter, code int, err error) {
	log.Println(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(map[string]string{"message": err.Error()}); err != nil {
		panic(err)
	}
}

func parseJson(r *http.Request, a interface{}) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	return decoder.Decode(a)
}

func queryStringToInt(qs url.Values, name string) (int, error) {
	if len(qs[name]) == 0 {
		return -1, nil
	}

	return strconv.Atoi(qs[name][0])
}
