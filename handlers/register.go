package handlers

import (
	"fmt"
	"github.com/sr8e/mellow-ir/auth"
	"log"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	url, state := auth.GenerateAuthURL()
	cookie := http.Cookie{
		Name:  "state",
		Value: state,
	}
	ce, err := auth.NewCookieEncrypter()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "invalid server configuration")
		log.Print(err)
		return
	}
	encrypted, err := ce.Encode(cookie)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "something is wrong...")
		log.Print(err)
		return
	}
	http.SetCookie(w, &encrypted)

	fmt.Fprintf(w, `<a href="%s">Log in with discord</a>`, url)
}
