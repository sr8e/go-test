package handlers

import (
	"github.com/sr8e/mellow-ir/auth"
	"fmt"
	"net/http"
)
func Callback(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("state")
	query := r.URL.Query()
	queryState := query.Get("state")
	if err != nil || c.Value != queryState {
		fmt.Fprintf(w, "authentication failed: invalid state")
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
	fmt.Fprintf(w, `authentication succeed: %s, %s, <img src="https://cdn.discordapp.com/avatar/%[1]s/%[3]s.png">`, user.Id, user.UserName, user.Avatar)
}
