package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func newAPIRouter() http.Handler {
	router := mux.NewRouter()
	chain := alice.New(stripPrefixMiddleware("/api"))

	for _, route := range apiRoutes {
		handlerFunc := route.handlerFunc

		// add middlleware here

		router.
			Methods(route.method).
			Path(route.pattern).
			Name(route.name).
			HandlerFunc(handlerFunc)

	}

	return chain.Then(router)
}
