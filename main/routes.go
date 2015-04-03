package main

type route struct {
	name    string
	method  string
	pattern string
	handler appHandler
}

type routes []route

var apiRoutes = routes{
	route{
		"NewLicense",
		"POST",
		"/licenses",
		NewLicense,
	},
	route{
		"RevokeLicense",
		"POST",
		"/licenses/{id}/revoke",
		RevokeLicense,
	},
}
