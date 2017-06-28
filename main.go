package main

import (
	"log"
	"net/http"

	"github.com/AlexanderSuv/go-blog-BE/router"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", router.New()))
}
