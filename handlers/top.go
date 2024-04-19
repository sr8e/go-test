package handlers

import (
	"fmt"
	"net/http"
)

func Top(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to Mell0wIR")
}
