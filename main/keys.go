package main

import (
	"crypto/rsa"

	"github.com/danielchatfield/go-jwt"
	"golang.org/x/net/context"
)

func getPrivateKey(c context.Context, kid string) (*rsa.PrivateKey, error) {
	file, err := getKey(c, kid, "private.pem")

	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPrivateKeyFromPEM(file)
}

func getPublicKey(c context.Context, kid string) (*rsa.PublicKey, error) {
	file, err := getKey(c, kid, "public.pem")

	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(file)
}

func getKey(c context.Context, kid string, fileName string) (key []byte, err error) {
	sc := NewStorageContext(c)

	// directory traversal is not a problem since this input is internal

	return sc.ReadFile("keys/" + kid + "/" + fileName)
}
