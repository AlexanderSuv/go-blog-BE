package main

import (
	"github.com/AlexanderSuv/go-blog-BE/router"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", router.New()))
}
