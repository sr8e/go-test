package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sr8e/mellow-ir/db"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := r.PostFormValue("id")
	pw := r.PostFormValue("password")

	dbUser := db.User{Id: id}
	ok, err := dbUser.VerifySecretToken(pw)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error on login: %s", err)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "ok")
}
