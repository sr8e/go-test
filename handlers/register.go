package handlers

import (
	"github.com/sr8e/mellow-ir/auth"
	"fmt"
	"net/http"
	"log"
)

func Register(w http.ResponseWriter, r *http.Request) {
	url, state := auth.GenerateAuthURL()
	cookie := http.Cookie{
		Name: "state",
		Value: state,
	}
	http.SetCookie(w, &cookie)

	_, err := fmt.Fprintf(w, `<a href="%s">Log in with discord</a>`, url)

	if err != nil {
		log.Print("error during handle register request...")
	}
}