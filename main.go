package main

import (
	"net/http"

	"github.com/sr8e/mellow-ir/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/login", http.HandlerFunc(handlers.Login))
	mux.Handle("/register", http.HandlerFunc(handlers.Register))
	mux.Handle("/callback", http.HandlerFunc(handlers.Callback))
	http.ListenAndServe(":8080", mux)
}
