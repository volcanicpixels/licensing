package main

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Result interface{} `json:"result"`
}

func writeJSON(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(response{v})
}
