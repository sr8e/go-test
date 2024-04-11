package auth

import (
	"os"
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

func GenerateAuthURL() (authUrl string, state string) {
	// generate state string
	randbytes := make([]byte, 16)
	rand.Read(randbytes)
	state = hex.EncodeToString(randbytes)

	query := url.Values{}
	query.Add("client_id", os.Getenv("DISCORD_CLIENT_ID"))
	query.Add("redirect_uri", os.Getenv("CALLBACK_URL"))
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
	postBody.Add("redirect_uri", os.Getenv("CALLBACK_URL"))

	req, err := http.NewRequest(http.MethodPost, "https://discord.com/api/oauth2/token", strings.NewReader(postBody.Encode()))
	if err != nil {
		log.Print("could not construct request")
		return
	}
	
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"))

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
