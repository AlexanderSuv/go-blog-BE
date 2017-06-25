package handlers

import (
	"net/http"
)

type RestAPI struct {
	GET    func(http.ResponseWriter, *http.Request)
	POST   func(http.ResponseWriter, *http.Request)
	PUT    func(http.ResponseWriter, *http.Request)
	DELETE func(http.ResponseWriter, *http.Request)
}

func (rest *RestAPI) HandleHTTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			rest.GET(w, r)
		case "POST":
			rest.POST(w, r)
		case "PUT":
			rest.PUT(w, r)
		case "DELETE":
			rest.DELETE(w, r)
		}
	}
}
