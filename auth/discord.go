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
	"time"
	"github.com/sr8e/mellow-ir/db"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn int `json:"expires_in"`
	Expire time.Time
}

type DiscordUser struct {
	Id string
	UserName string
	Avatar string
}

func wrapClient(r *http.Request) (*http.Response, []byte, error) {
	client := http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		err = fmt.Errorf("could not send request: %w", err)
		return resp, nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("could not read response: %w", err)
		return resp, body, err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("request failed: %s, content: %s", resp.Status, string(body))
		return resp, body, err
	}

	return resp, body, nil
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

func GetAuthToken(code string, refresh bool) (tr TokenResponse, err error) {
	postBody := url.Values{}
	if refresh {
		postBody.Add("grant_type", "refresh_token")
		postBody.Add("refresh_token", code)
	} else {
		postBody.Add("grant_type", "authorization_code")
		postBody.Add("code", code)
		postBody.Add("redirect_uri", discordCallbackURL)
	}

	req, err := http.NewRequest(http.MethodPost, "https://discord.com/api/oauth2/token", strings.NewReader(postBody.Encode()))
	if err != nil {
		err = fmt.Errorf("could not construct request: %w", err)
		return
	}
	
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(discordClientId, discordClientSecret)

	resp, body, err := wrapClient(req)
	if err != nil {
		return
	}

	respTime, err := time.Parse(time.RFC1123, resp.Header.Get("Date"))
	if err != nil {
		log.Printf("cannot parse date (%s), set local time instead")
		respTime = time.Now()
	}

	err = json.Unmarshal(body, &tr)
	if err != nil {
		err = fmt.Errorf("could not unmarshal json: %w", err)
		return
	}

	tr.Expire = respTime.Add(time.Second * time.Duration(tr.ExpiresIn))
	return tr, nil
}

func GetUser(token string) (user DiscordUser, err error) {
	req, err := http.NewRequest(http.MethodGet, "https://discord.com/api/oauth2/@me", nil)
	if err != nil {
		err = fmt.Errorf("could not construct request: %w", err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	_, body, err := wrapClient(req)
	if err != nil {
		return
	}

	type UserResponse struct {
		User DiscordUser
	}
	var ur UserResponse
	err = json.Unmarshal(body, &ur)
	if err != nil {
		err = fmt.Errorf("could not unmarshal json: %w", err)
		return
	}
	return ur.User, nil
}

func FromDiscordUser(u *db.User, tr TokenResponse, du DiscordUser) {
	u.DisplayName = du.UserName
	u.IconURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", du.Id, du.Avatar)
	u.AccessToken = tr.AccessToken
	u.RefreshToken = tr.RefreshToken
	u.Expire = tr.Expire
}
