package handlers

import (
	"net/http"
)



func Authors() http.HandlerFunc {
	handler := RestAPI{}

	return handler.HandleHTTP()
}
