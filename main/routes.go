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
	route{
		"DecodeLicense",
		"POST",
		"/licenses/_/decode",
		DecodeLicense,
	},
	route{
		"UpdateRevocationFile",
		"GET",
		"/update_revocation_file",
		UpdateRevocationFile,
	},
}
