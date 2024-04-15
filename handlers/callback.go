package handlers

import (
	"github.com/sr8e/mellow-ir/auth"
	"fmt"
	"log"
	"time"
	"net/http"
)
func Callback(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("state")
	if err != nil {
		fmt.Fprintf(w, "authentication failed: cookie not exist")
		return
	}	

	ce, err := auth.NewCookieEncrypter()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "invalid server configuration")
		log.Print(err)
		return
	}
	value, err := ce.Decode(*c)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "authentication failed: cookie corrupted")
		log.Print(err)
		return
	}
	query := r.URL.Query()
	queryState := query.Get("state")
	if value != queryState {
		w.WriteHeader(400)
		fmt.Fprintf(w, "authentication failed: invalid state")
		log.Printf("cookie:%s, query:%s", value, queryState)
		return
	}

	token, err := auth.GetAuthToken(query.Get("code"))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "authentication failed: could not acquire token")
		log.Printf("could not acquire token: %w", err)
		return
	}

	user, err := auth.GetUser(token.AccessToken)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "could not fetch user data")
		log.Printf("could not fetch user data: %w", err)
		return
	}

	// set session cookie
	sessRaw := http.Cookie{
		Name: "session",
		Value: user.Id,
		Expires: time.Now().AddDate(0, 0, 7),
	}
	sessCookie, err := ce.Encode(sessRaw)
	if err != nil {
		log.Print("session cookie encryption failed")
	} else {
		http.SetCookie(w, &sessCookie)
	}

	// delete state cookie
	c.MaxAge = -1
	http.SetCookie(w, c)

	fmt.Fprintf(w, `authentication succeed: %s, %s, <img src="https://cdn.discordapp.com/avatar/%[1]s/%[3]s.png">`, user.Id, user.UserName, user.Avatar)
}
