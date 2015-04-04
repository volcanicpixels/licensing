package main

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/danielchatfield/go-jwt"
	"github.com/gorilla/mux"
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

// NewLicense handles POST requests on /api/licenses/create
//
// The request body must contain a JSON object with a product field
//
// Examples:
//
//  POST /api/licenses {"product": ""}
//  400 empty title
//
//  POST /api/licenses {"product": "domain_changer"}
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

func revokeLicense(c context.Context, id string) error {
	// ideally we would simply add the license ID on to the end of the revocations.txt file
	// but Google Storage doesn't support appends.
	// It does support a composition operation, so we could write the new ID to a new file
	// and then compose the original with the new one to ensure atomicity, except the Google
	// storage client library does not implement this operation.
	// Therefore the best we can do without stupidly complex locks is to simply read in the current file
	// and then write a new file with the addition

	// read the current revocations.txt file
	sc := NewStorageContext(c)
	data, err := sc.ReadFile("revocations.txt")

	if err != nil {
		return err
	}

	line := id

	// almost certainly a better way to do this
	data = []byte(string(data) + "\n" + line)

	if err := sc.WriteFile("revocations.txt", data); err != nil {
		return err
	}

	return nil
}

// RevokeLicense handles POST requests to /api/licenses/{ID}/revoke
func RevokeLicense(c context.Context, w http.ResponseWriter, r *http.Request) *appError {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := revokeLicense(c, id); err != nil {
		return &appError{err, "An error occurred updating the revocations file", http.StatusInternalServerError}
	}

	writeJSON(w, 200, "SUCCESS")

	return nil
}

func DecodeLicense(c context.Context, w http.ResponseWriter, r *http.Request) *appError {
	var req struct{ License string }
	var err error

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &appError{err, "Could not decode json request", http.StatusBadRequest}
	}

	// req successfully Decoded

	var key *rsa.PublicKey
	if key, err = getPublicKey(c, "plugin"); err != nil {
		return &appError{err, "Could not load public key for verifying", http.StatusInternalServerError}
	}

	l, err := license.Parse(req.License, key)

	if err != nil {
		return &appError{err, "An error occured parsing the token", http.StatusBadRequest}
	}

	// license successfuly decoded - now lets return the response
	writeJSON(w, 200, l)

	return nil
}

func UpdateRevocationFile(c context.Context, w http.ResponseWriter, r *http.Request) *appError {
	sc := NewStorageContext(c)
	data, err := sc.ReadFile("revocations.txt")

	if err != nil {
		return &appError{err, "An error occurred reading the revocations.txt file", http.StatusInternalServerError}
	}

	revocations := strings.Split(string(data), "\n")

	var formatted []string

	for _, line := range revocations {
		id := strings.TrimSpace(strings.Split(line, "#")[0])

		if id == "" {
			continue
		}

		formatted = append(formatted, id)
	}

	// ok, we have the revocations now

	t := jwt.NewToken(jwt.RSA)
	t.SetClaim("_revoked", formatted)
	t.SetClaim("exp", time.Now().Add(time.Hour*72).Unix())

	key, err := getPrivateKey(c, "plugin")

	if err != nil {
		return &appError{err, "The private key could not be retrieved", http.StatusInternalServerError}
	}

	type revokedJSON struct {
		Token string `json:"token"`
	}

	tokenString, err := t.Encode(key)

	if err != nil {
		return &appError{err, "An error occured when signing the token", http.StatusInternalServerError}
	}

	rev, err := json.Marshal(revokedJSON{tokenString})

	if err != nil {
		return &appError{err, "An error occured when marshalling the JSON", http.StatusInternalServerError}
	}

	err = sc.WriteFile("revocations.json", []byte(rev))

	if err != nil {
		return &appError{err, "An error occured when writing the revocations.json file", http.StatusInternalServerError}
	}

	err = sc.MakePublic("revocations.json")

	if err != nil {
		return &appError{err, "An error occured when making the revocations.json file public", http.StatusInternalServerError}
	}

	writeJSON(w, 200, "SUCCESS")

	return nil

}
