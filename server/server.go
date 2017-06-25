package server

import (
	"net/http"

	"github.com/AlexanderSuv/goblog/server/handlers"
)

func Start() {
	h := http.NewServeMux()

	h.Handle("/authors/", handlers.Authors())
	http.ListenAndServe(":8080", h)
}
