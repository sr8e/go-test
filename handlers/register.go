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
	ce := auth.NewCookieEncrypter()
	encrypted, err := ce.Encode(cookie)
	if err != nil {
		fmt.Fprint(w, "something is wrong...")
		log.Print(err)
		return
	}
	http.SetCookie(w, &encrypted)

	fmt.Fprintf(w, `<a href="%s">Log in with discord</a>`, url)

	if err != nil {
		log.Print("error during handle register request...")
	}
}
