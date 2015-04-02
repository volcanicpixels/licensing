package main

import "net/http"

func stripPrefixMiddleware(prefix string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.StripPrefix(prefix, h)
	}
}
