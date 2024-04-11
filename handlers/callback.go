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

	ce := auth.NewCookieEncrypter()
	value, err := ce.Decode(*c)
	if err != nil {
		fmt.Fprintf(w, "authentication failed: cookie corrupted")
		log.Print(err)
	}
	query := r.URL.Query()
	queryState := query.Get("state")
	if value != queryState {
		fmt.Fprintf(w, "authentication failed: invalid state")
		log.Printf("cookie:%s, query:%s", value, queryState)
		return
	}

	token := auth.GetAuthToken(query.Get("code"))
	if token == "" {
		fmt.Fprintf(w, "authentication failed: could not acquire token")
		return
	}

	user, err := auth.GetUser(token)
	if err != nil {
		fmt.Fprintf(w, "could not fetch user data")
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
