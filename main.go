package main

import (
	"github.com/sr8e/mellow-ir/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/login", http.HandlerFunc(handlers.Login))
	mux.Handle("/register", http.HandlerFunc(handlers.Register))
	mux.Handle("/callback", http.HandlerFunc(handlers.Callback))
	http.ListenAndServe(":8080", mux)
}
