package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func newAPIRouter() http.Handler {
	router := mux.NewRouter()
	chain := alice.New(stripPrefixMiddleware("/api"))

	for _, route := range apiRoutes {
		handler := appHandler(route.handler)

		// add middlleware here

		router.
			Methods(route.method).
			Path(route.pattern).
			Name(route.name).
			Handler(handler)

	}

	//router.HandleFunc("/licenses", testHandler)

	return chain.Then(router)
}
func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hi there")
}
