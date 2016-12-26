package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := NewResponse(w)
		response.Body.Error = true
		response.Body.Message = fmt.Sprintf("%s not found.", r.RequestURI)
		response.Status = http.StatusNotFound
		response.Write()
	})

	return router
}

var routes = Routes{
	Route{
		Name: "Index",
		Method: "POST",
		Pattern: "/",
		HandlerFunc: Index,
	},
}
