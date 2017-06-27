package router

import (
	"net/http"

	"github.com/AlexanderSuv/go-blog-BE/controllers"
	"github.com/gorilla/mux"
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
	addRoute(&Route{"GET", "/authors", controllerAuthors.Get}, apiRouter)
	addRoute(&Route{"POST", "/authors", controllerAuthors.Post}, apiRouter)
	addRoute(&Route{"GET", "/authors/{id}", controllerAuthors.GetById}, apiRouter)
	addRoute(&Route{"PUT", "/authors/{id}", controllerAuthors.Put}, apiRouter)

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
