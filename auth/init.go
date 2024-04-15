package auth

import (
	"os"
	"log"
	"encoding/base64"
)

var discordClientId string
var discordClientSecret string
var discordCallbackURL string

var secretKey []byte

func init() {
	discordClientId = os.Getenv("DISCORD_CLIENT_ID")
	if discordClientId == "" {
		log.Println("environment variable DISCORD_CLIENT_ID is not set")
	}

	discordClientSecret = os.Getenv("DISCORD_CLIENT_SECRET")
	if discordClientSecret == "" {
		log.Println("environment variable DISCORD_CLIENT_SECRET is not set")
	}

	discordCallbackURL = os.Getenv("CALLBACK_URL")
	if discordCallbackURL == "" {
		log.Println("environment variable CALLBACK_URL is not set")
	}

	secKeyStr := os.Getenv("SECRET_KEY")
	if secKeyStr == "" {
		log.Println("envitonment variable SECRET_KEY is not set")
	} else {
		// try base64 decode
		b, err := base64.StdEncoding.DecodeString(secKeyStr)
		if err != nil {
			log.Println("could not decode SECRET_KEY. directly convert to byte array")
			b = []byte(secKeyStr)
		}
		secretKey = make([]byte, 32)
		copy(secretKey, b[:min(len(b), 32)])
	}
}
