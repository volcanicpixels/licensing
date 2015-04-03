package main

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/volcanicpixels/licensing/license"
)

type appHandler func(context.Context, http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error   error
	Message string
	Code    int
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if e := fn(c, w, r); e != nil {
		log.Errorf(c, "[%v] %v", e.Message, e.Error)
		http.Error(w, e.Message, e.Code)
	}
}

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
func NewLicense(c context.Context, w http.ResponseWriter, r *http.Request) *appError {
	var req struct{ Product string }
	var err error

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &appError{err, "Could not decode json request", http.StatusBadRequest}
	}

	var key *rsa.PrivateKey
	if key, err = getPrivateKey(c, "plugin"); err != nil {
		return &appError{err, "Could not load private key for signing", http.StatusInternalServerError}
	}

	// create the license
	lic := license.New(req.Product)

	var licStr string
	if licStr, err = lic.Encode(key); err != nil {
		return &appError{err, "Could not encode the license", http.StatusInternalServerError}
	}

	writeJSON(w, 200, licStr)
	return nil
}
