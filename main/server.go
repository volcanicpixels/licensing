package main

import "net/http"

func init() {
	http.Handle("/api/", newAPIRouter())
}
