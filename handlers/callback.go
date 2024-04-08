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
	fmt.Println(token)
	fmt.Fprint(w, "authentication succeed")
}