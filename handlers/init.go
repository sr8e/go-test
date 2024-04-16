package handlers

import (
	"log"
	"os"
)

var cryptKey string

func init() {
	cryptKey = os.Getenv("IR_CRYPT_KEY")

	if cryptKey == "" {
		log.Print("environment variable IR_CRYPT_KEY is empty")
	}
}
