package handlers

import (
	"fmt"
	"net/http"
	"log"
)

func Login(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "cryptkey=%s", cryptKey)

	if err != nil {
		log.Print("error during handle login request...")
	}
}