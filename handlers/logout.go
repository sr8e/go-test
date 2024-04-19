package handlers

import (
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	defer func() {
		http.Redirect(w, r, "/", http.StatusFound)
	}()

	c, _ := r.Cookie("session")
	if c == nil {
		return
	}
	c.MaxAge = -1
	http.SetCookie(w, c)
}
