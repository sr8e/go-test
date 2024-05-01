package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

func Encrypt(value string, secretKey []byte) (string, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	iv := make([]byte, bs)
	rand.Read(iv)

	valByte := []byte(value)
	cipher.NewCTR(block, iv).XORKeyStream(valByte, valByte)

	return base64.StdEncoding.EncodeToString(append(iv, valByte...)), nil
}

func Decrypt(cryptStr string, secretKey []byte) (decrypted string, err error) {
	cryptByte, err := base64.StdEncoding.DecodeString(cryptStr)
	if err != nil {
		err = fmt.Errorf("could not decode crypted body: %w", err)
		return
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return
	}
	bs := block.BlockSize()
	if len(cryptByte) <= bs {
		err = errors.New("crypted content too short")
		return
	}
	iv := cryptByte[:bs]
	body := cryptByte[bs:]

	cipher.NewCTR(block, iv).XORKeyStream(body, body)

	return string(body), nil
}
