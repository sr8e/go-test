package auth

import (
	"os"
	"fmt"
	"io"
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

func GetAuthToken(code string) string {
	postBody := url.Values{}
	postBody.Add("grant_type", "authorization_code")
	postBody.Add("code", code)
	postBody.Add("redirect_uri", os.Getenv("CALLBACK_URL"))

	resp, err := http.PostForm("https://discord.com/api/oauth2/token", postBody)
	if err != nil {
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	fmt.Println(string(body))

	var obj map[string] any
	err = json.Unmarshal(body, obj)
	if err != nil {
		return ""
	}
	return obj["access_token"].(string)
}
