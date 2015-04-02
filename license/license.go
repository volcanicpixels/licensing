package license

import (
	"time"

	"github.com/danielchatfield/go-jwt"
	"github.com/dchest/uniuri"
)

// License represents a license entity
type License interface {
	Encode(key interface{}) (string, error)
}

type license struct {
	id        string
	product   string
	createdAt time.Time
	attrs     map[string]interface{}
}

// New creates a new License. Takes the product that the license is for.
func New(product string) License {
	return &license{
		id:        uniuri.New(),
		product:   product,
		createdAt: time.Now(),
		attrs:     make(map[string]interface{}),
	}
}

func (l *license) Encode(key interface{}) (encoded string, err error) {
	// This is essentially a wrapper around jwt.Token.Encode()

	t := jwt.NewToken(jwt.RSA)

	t.SetClaim("jti", l.id)
	t.SetClaim("iat", l.createdAt.Unix())
	t.SetClaim("_prod", l.product)
	t.SetClaim("_attrs", l.attrs)

	return t.Encode(key)
}
