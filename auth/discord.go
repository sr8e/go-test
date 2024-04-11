package auth

import (
	"fmt"
	"io"
	"log"
	"strings"
	"net/url"
	"net/http"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
)

type DiscordUser struct {
	Id string
	UserName string
	Avatar string
}

func GenerateAuthURL() (authUrl string, state string) {
	// generate state string
	randbytes := make([]byte, 16)
	rand.Read(randbytes)
	state = hex.EncodeToString(randbytes)

	query := url.Values{}
	query.Add("client_id", discordClientId)
	query.Add("redirect_uri", discordCallbackURL)
	query.Add("response_type", "code")
	query.Add("scope", "identify")
	query.Add("state", state)

	authUrl = fmt.Sprintf("https://discord.com/oauth2/authorize?%s", query.Encode())

	return authUrl, state
}

func GetAuthToken(code string) (token string) {
	postBody := url.Values{}
	postBody.Add("grant_type", "authorization_code")
	postBody.Add("code", code)
	postBody.Add("redirect_uri", discordCallbackURL)

	req, err := http.NewRequest(http.MethodPost, "https://discord.com/api/oauth2/token", strings.NewReader(postBody.Encode()))
	if err != nil {
		log.Print("could not construct request")
		return
	}
	
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(discordClientId, discordClientSecret)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("could not send request")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("could not read response")
		return
	}
	
	if resp.StatusCode != 200 {
		log.Printf("request failed: %s, content: %s", resp.Status, string(body))
		return
	}

	type TokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	var tr TokenResponse
	err = json.Unmarshal(body, &tr)
	if err != nil {
		log.Printf("could not unmarshal json")
		return
	}
	return tr.AccessToken
}

func GetUser(token string) (user DiscordUser, err error) {
	req, err := http.NewRequest(http.MethodGet, "https://discord.com/api/oauth2/@me", nil)
	if err != nil {
		log.Printf("could not construct request")
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("could not send request")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("could not read response")
		return
	}

	if resp.StatusCode != 200 {
		log.Printf("request failed: %s, content: %s", resp.Status, resp.Body)
		return
	}
	
	type UserResponse struct {
		User DiscordUser
	}
	var ur UserResponse
	err = json.Unmarshal(body, &ur)
	if err != nil {
		log.Printf("could not unmarshal json")
		return
	}
	return ur.User, nil
}
