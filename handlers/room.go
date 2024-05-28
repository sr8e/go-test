package handlers

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/sr8e/mellow-ir/db"
)

func CreateRoom(req chan<- string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("incoming request")
		_, ok, err := BearerAuth(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error at api/room/create: %s", err)
			return
		}
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ut := uint64(time.Now().UnixMilli())
		randn := rand.Uint64()
		b := make([]byte, 0)
		b = binary.BigEndian.AppendUint64(b, ut)
		b = binary.BigEndian.AppendUint64(b, randn)
		roomName := hex.EncodeToString(b)
		roomKey, err := db.SetRoomKey(roomName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error on creating roomkey: %s", err)
			return
		}
		req <- roomName
		fmt.Fprintf(w, "OK:%s,%s", roomName, roomKey)
	}
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	kw := r.FormValue("keyword")
	if kw == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	roomName, ok, err := db.GetRoomName(kw)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error on finding room: %s", err)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "OK:%s", roomName)
}
