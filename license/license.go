package license

import (
	"errors"
	"time"

	"github.com/danielchatfield/go-jwt"
	"github.com/dchest/uniuri"
)

// License represents a license entity
type License interface {
	Encode(key interface{}) (string, error)
}

type license struct {
	ID       string                 `json:"id"`
	Product  string                 `json:"product"`
	IssuedAt time.Time              `json:"issuedAt"`
	Attrs    map[string]interface{} `json:"attrs"`
}

// New creates a new License. Takes the product that the license is for.
func New(product string) License {
	return &license{
		ID:       uniuri.New(),
		Product:  product,
		IssuedAt: time.Now(),
		Attrs:    make(map[string]interface{}),
	}
}

func (l *license) Encode(key interface{}) (encoded string, err error) {
	// This is essentially a wrapper around jwt.Token.Encode()

	t := jwt.NewToken(jwt.RSA)

	t.SetClaim("jti", l.ID)
	t.SetClaim("iat", l.IssuedAt.Unix())
	t.SetClaim("_prod", l.Product)
	t.SetClaim("_attrs", l.Attrs)

	return t.Encode(key)
}

func Parse(token string, key interface{}) (*license, error) {
	tok, err := jwt.ParseToken(token, jwt.RSA, key)

	if err != nil {
		return nil, err
	}

	// tok decoded - now lets construct the license

	l := &license{}
	var ok bool

	l.ID, ok = tok.Claim("jti").(string)

	if !ok {
		return nil, errors.New("Error extracting license ID")
	}

	l.Product, ok = tok.Claim("_prod").(string)

	if !ok {
		return nil, errors.New("Error extracting license product")
	}

	// these fields are not strictly required so we don't handle errors
	timestamp := int64(tok.Claim("iat").(float64))
	l.IssuedAt = time.Unix(timestamp, 0)
	l.Attrs = tok.Claim("_attrs").(map[string]interface{})

	return l, nil
}
