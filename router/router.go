package router

import (
	"github.com/AlexanderSuv/go-blog-BE/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

var API_VERSION_PREFIX = "/api/v1/"

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func New() *mux.Router {
	apiRouter := mux.NewRouter().StrictSlash(true)

	controllerAuthors := controllers.Authors{}
	addRoute(&Route{"GET", "/authors/", controllerAuthors.Get}, apiRouter)
	addRoute(&Route{"POST", "/authors", controllerAuthors.Post}, apiRouter)
	addRoute(&Route{"GET", "/authors/{id}", controllerAuthors.GetById}, apiRouter)
	addRoute(&Route{"POST", "/authors/{id}", controllerAuthors.Post}, apiRouter)

	return apiRouter
}

func addRoute(route *Route, router *mux.Router) {
	handler := Logger(route.HandlerFunc)
	router.
		PathPrefix(API_VERSION_PREFIX).
		Methods(route.Method).
		Path(route.Pattern).
		Handler(handler)
}
