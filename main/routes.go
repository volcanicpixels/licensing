package main

import "net/http"

type route struct {
	name        string
	method      string
	pattern     string
	handlerFunc http.HandlerFunc
}

type routes []route

var apiRoutes = routes{
	route{
		"NewLicense",
		"POST",
		"/licenses",
		NewLicense,
	},
}
