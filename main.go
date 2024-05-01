package main

import (
	"context"
	"net/http"

	"github.com/sr8e/mellow-ir/battle"
	"github.com/sr8e/mellow-ir/handlers"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	newRoomChan := make(chan string)
	defer close(newRoomChan)

	go battle.BattleMainRoutine(ctx, cancel, newRoomChan)

	mux := http.NewServeMux()
	// pages
	mux.Handle("/", http.HandlerFunc(handlers.Top))
	mux.Handle("/mypage", handlers.RequireSession(handlers.MyPage))
	mux.Handle("/authenticate", handlers.WithoutSession(handlers.Register))
	mux.Handle("/callback", http.HandlerFunc(handlers.Callback))
	mux.Handle("/logout", http.HandlerFunc(handlers.Logout))
	// apis
	mux.Handle("/api/me", http.HandlerFunc(handlers.Login))
	mux.Handle("/api/room/create", http.HandlerFunc(handlers.CreateRoom(newRoomChan)))
	http.ListenAndServe(":8080", mux)
}
