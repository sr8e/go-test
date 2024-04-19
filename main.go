package main

import (
	"net/http"

	"github.com/sr8e/mellow-ir/handlers"
)

func main() {
	mux := http.NewServeMux()
	// pages
	mux.Handle("/", http.HandlerFunc(handlers.Top))
	mux.Handle("/mypage", handlers.RequireSession(handlers.MyPage))
	mux.Handle("/authenticate", handlers.WithoutSession(handlers.Register))
	mux.Handle("/callback", http.HandlerFunc(handlers.Callback))
	mux.Handle("/logout", http.HandlerFunc(handlers.Logout))
	// apis
	mux.Handle("/api/login", http.HandlerFunc(handlers.Login))
	http.ListenAndServe(":8080", mux)
}
