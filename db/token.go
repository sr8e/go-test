package db

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func (u *User) GenerateSecretToken() (secret string) {
	userSecret := make([]byte, 16)
	secretSalt := make([]byte, 16)
	rand.Read(userSecret)
	rand.Read(secretSalt)
	secretHash := sha256.Sum256(append(userSecret, secretSalt...))
	u.secretHash = hex.EncodeToString(secretHash[:])
	u.secretSalt = hex.EncodeToString(secretSalt)
	return hex.EncodeToString(userSecret)
}

func (u *User) VerifySecretToken(secret string) (bool, error) {
	if !u.loaded {
		ok, err := u.Get()
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	secretBytes, err := hex.DecodeString(secret)
	if err != nil { // invalid string
		return false, nil
	}
	saltBytes, _ := hex.DecodeString(u.secretSalt)
	hash := sha256.Sum256(append(secretBytes, saltBytes...))

	return hex.EncodeToString(hash[:]) == u.secretHash, nil
}
