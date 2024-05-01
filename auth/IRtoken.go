package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/sr8e/mellow-ir/db"
)

func VerifyIRToken(token string) (string, bool, error) {
	decToken, err := Decrypt(token, secretKey)
	log.Printf("dectoken: %s", decToken)
	if err != nil {
		return "", false, err
	}
	id, body, ok := strings.Cut(decToken, ":")
	if !ok {
		return "", false, nil
	}
	stored, ok, err := db.ReadIRSession(id)
	if err != nil {
		return "", false, err
	}
	return id, ok && (body == stored), nil
}

func CreateIRToken(id string) (string, error) {
	token, stored, _ := db.ReadIRSession(id)
	if stored {
		err := db.RefreshIRSession(id)
		if err != nil {
			return "", err
		}
	} else {

		randByte := make([]byte, 16)
		rand.Read(randByte)
		token = base64.StdEncoding.EncodeToString(randByte)
		err := db.WriteIRSession(id, token)
		if err != nil {
			return "", err
		}
	}
	encToken, err := Encrypt(fmt.Sprintf("%s:%s", id, token), secretKey)
	if err != nil {
		return "", err
	}

	return encToken, nil
}
