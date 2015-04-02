package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func init() {

	apiChain := alice.New(stripPrefixMiddleware("/api"))
	apiRouter := mux.NewRouter()

	apiRouter.HandleFunc("/licenses", NewLicense).Methods("POST")

	http.Handle("/api/", apiChain.Then(apiRouter))
}
