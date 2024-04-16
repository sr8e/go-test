package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sr8e/mellow-ir/auth"
	"github.com/sr8e/mellow-ir/db"
)

func RequireSession(mainHandler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	preHandler := func(w http.ResponseWriter, r *http.Request) {
		u, err := checkSession(r)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal server error")
			log.Printf("session check failed: %s", err)
			return
		}
		t := time.Now()
		if u == nil || t.After(u.Expire) {
			http.Redirect(w, r, "/login", 301)
			return
		}
		if u.Expire.Sub(t).Seconds() < 3600*24*3.5 {
			// half of ttl, try refresh access token async
			go func() {
				tr, err := auth.GetAuthToken(u.RefreshToken, true)
				if err != nil {
					log.Printf("token refresh failed: %s", err)
					return
				}
				auth.FromDiscordUser(u, &tr, nil)
				err = u.Save()
				if err != nil {
					log.Printf("user update failed: %s", err)
					return
				}
			}()
		}
		mainHandler(w, r)
	}
	return http.HandlerFunc(preHandler)
}

func checkSession(r *http.Request) (*db.User, error) {
	sess, err := r.Cookie("session")
	if err != nil { // cookie does not exist
		log.Print("no cookie")
		return nil, nil
	}
	ce, err := auth.NewCookieEncrypter()
	if err != nil {
		return nil, err
	}
	id, err := ce.Decode(*sess)
	if err != nil { // invalid cookie
		log.Printf("cookie decrypt fail %s", err)
		return nil, nil
	}
	u := db.User{Id: id}
	ok, err := u.Get()
	if err != nil {
		return nil, err
	}
	if !ok {
		log.Print("user not found")
		return nil, nil
	}
	return &u, nil
}
