package auth

import (
	"github.com/sr8e/mellow-ir/db"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateSecretToken() (secret, hash, salt string){
	userSecret := make([]byte, 16)
	secretSalt := make([]byte, 16)
	rand.Read(userSecret)
	rand.Read(secretSalt)
	secretHash := sha256.Sum256(append(userSecret, secretSalt...))
	return hex.EncodeToString(userSecret), hex.EncodeToString(secretHash[:]), hex.EncodeToString(secretSalt)
}

func VerifySecretToken(id, secret string) (bool, error) {
	dbUser := db.User{Id: id}
	ok, err := dbUser.Get()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	secretBytes, err := hex.DecodeString(secret)
	if err != nil { // invalid string
		return false, nil
	}
	saltBytes, _ := hex.DecodeString(dbUser.SecretSalt)
	hash := sha256.Sum256(append(secretBytes, saltBytes...))

	return hex.EncodeToString(hash[:]) == dbUser.SecretHash, nil
}
