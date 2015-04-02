package main

import (
	"encoding/json"
	"net/http"

	"github.com/volcanicpixels/licensing/license"
)

// NewLicense handles POST requests on /api/licenses
//
// The request body must contain a JSON object with a product field
//
// Examples:
//
//  POST /api/licenses {"product": ""}
//  400 empty title
//
//  POST /api/license {"product": "domain_changer"}
//  200
func NewLicense(w http.ResponseWriter, r *http.Request) {
	var req struct{ Product string }

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// create the license
	lic := license.New(req.Product)
	if licStr, err := lic.Encode(getPrivateKey(req.Product)); err == nil {
		writeJSON(w, 200, licStr)
		return
	}

	http.Error(w, "An error occured when generating the license", http.StatusInternalServerError)
}
