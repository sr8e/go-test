package auth

import (
	"os"
	"log"
)

var discordClientId string
var discordClientSecret string
var discordCallbackURL string

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
}
