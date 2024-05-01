package handlers

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
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
		req <- roomName
		fmt.Fprintf(w, "OK:%s", roomName)
	}
}
