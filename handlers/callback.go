package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sr8e/mellow-ir/auth"
	"github.com/sr8e/mellow-ir/db"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("state")
	if err != nil {
		w.WriteHeader(400)
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

	token, err := auth.GetAuthToken(query.Get("code"), false)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "authentication failed: could not acquire token")
		log.Printf("could not acquire token: %s", err)
		return
	}

	user, err := auth.GetUser(token.AccessToken)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "could not fetch user data")
		log.Printf("could not fetch user data: %s", err)
		return
	}

	// set session cookie
	sessRaw := http.Cookie{
		Name:    "session",
		Value:   user.Id,
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

	// save on db
	dbUser := db.User{Id: user.Id}
	ok, err := dbUser.Get()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "db operation error")
		log.Printf("error on querying user: %s", err)
		return
	}
	auth.FromDiscordUser(&dbUser, &token, &user)
	if ok {
		// existing user
		err = dbUser.Save()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "db operation error")
			log.Printf("error on saving user: %s", err)
			return
		}
		http.Redirect(w, r, "/", 301)
		return
	}

	// create new user
	secret, hash, salt := auth.GenerateSecretToken()
	dbUser.SecretHash = hash
	dbUser.SecretSalt = salt
	err = dbUser.Save()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "db operation error")
		log.Printf("error on creating user: %s", err)
		return
	}

	fmt.Fprintf(w, `authentication succeed.
		Your User ID: %s,
		Your Secret token: %s (make sure you keep this code. it will never be shown again!)`,
		dbUser.Id, secret,
	)
}
