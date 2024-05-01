package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sr8e/mellow-ir/auth"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	u, ok, err := BasicAuth(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error at api/me: %s", err)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	irToken, err := auth.CreateIRToken(u.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error at api/me: %s", err)
		return
	}

	fmt.Fprintf(w, "%s:%s", u.DisplayName, irToken)
}
