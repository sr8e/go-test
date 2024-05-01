package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/sr8e/mellow-ir/auth"
	"github.com/sr8e/mellow-ir/db"
)

func BasicAuth(r *http.Request) (*db.User, bool, error) {
	basic := r.Header.Get("Authorization")

	sigBody, ok := strings.CutPrefix(basic, "Basic ")
	if !ok {
		return nil, false, nil
	}
	sigStr, err := base64.StdEncoding.DecodeString(sigBody)
	if err != nil {
		// invalid string. suppress err
		return nil, false, nil
	}
	id, pw, ok := strings.Cut(string(sigStr), ":")
	if !ok {
		return nil, false, nil
	}

	u := &db.User{Id: id}
	ok, err = u.VerifySecretToken(pw)
	if !ok {
		return nil, false, err
	}
	return u, true, nil
}

func BearerAuth(r *http.Request) (*db.User, bool, error) {
	header := r.Header.Get("Authorization")
	log.Printf("header: %s", header)

	sigBody, ok := strings.CutPrefix(header, "Bearer ")
	if !ok {
		return nil, false, nil
	}
	id, ok, err := auth.VerifyIRToken(sigBody)
	if !ok {
		return nil, false, err
	}

	u := &db.User{Id: id}
	ok, err = u.Get()
	if !ok {
		return nil, false, err
	}
	return u, true, nil
}
