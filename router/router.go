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

	authorsController := controllers.Authors{}
	addRoute(&Route{"GET", "/authors", authorsController.Get}, apiRouter)
	addRoute(&Route{"POST", "/authors", authorsController.Post}, apiRouter)
	addRoute(&Route{"GET", "/authors/{id}", authorsController.GetById}, apiRouter)
	addRoute(&Route{"PUT", "/authors/{id}", authorsController.Put}, apiRouter)
	addRoute(&Route{"DELETE", "/authors/{id}", authorsController.Delete}, apiRouter)

	postsController := controllers.Posts{}
	addRoute(&Route{"GET", "/posts", postsController.Get}, apiRouter)
	addRoute(&Route{"POST", "/posts", postsController.Post}, apiRouter)
	addRoute(&Route{"GET", "/posts/{id}", postsController.GetById}, apiRouter)
	addRoute(&Route{"PUT", "/posts/{id}", postsController.Put}, apiRouter)
	addRoute(&Route{"DELETE", "/posts/{id}", postsController.Delete}, apiRouter)

	commentsController := controllers.Comments{}
	addRoute(&Route{"GET", "/posts/{postId}/comments", commentsController.Get}, apiRouter)
	addRoute(&Route{"POST", "/posts/{postId}/comments", commentsController.Post}, apiRouter)
	addRoute(&Route{"GET", "/posts/{postId}/comments/{commentId}", commentsController.GetById}, apiRouter)
	addRoute(&Route{"PUT", "/posts/{postId}/comments/{commentId}", commentsController.Put}, apiRouter)
	addRoute(&Route{"DELETE", "/posts/{postId}/comments/{commentId}", commentsController.Delete}, apiRouter)

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
