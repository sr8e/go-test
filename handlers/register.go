package handlers

import (
	"fmt"
	"net/http"
	"log"
)

func Register(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "%s...", "wooooosh")

	if err != nil {
		log.Print("error during handle register request...")
	}
}