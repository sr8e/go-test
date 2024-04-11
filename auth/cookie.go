package auth

import (
	"fmt"
	"log"
	"time"
	"errors"
	"bytes"
	"strconv"
	"net/http"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type CookieEncrypter struct {
	secretKey []byte
}

func NewCookieEncrypter() CookieEncrypter {
	return CookieEncrypter{
		secretKey: secretKey,
	}
}

func (ce *CookieEncrypter) Encode(c http.Cookie) (out http.Cookie, err error) {
	out = c
	// encrypt value
	if ce.secretKey == nil {
		err = errors.New("secret key not set")
		return
	}
	enc, err := encrypt(c.Value, ce.secretKey)
	if err != nil {
		return
	} 
	
	encoded := signature(enc, c.Expires, ce.secretKey)
	out.Value = encoded

	err = out.Valid()
	if err != nil {
		return
	}
	return out, nil
}

func (ce *CookieEncrypter) Decode(c http.Cookie) (value string, err error) {
	b, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		return
	}
	
	cryptStr, err := verify(b, ce.secretKey)
	if err != nil {
		return
	}
	return decrypt(cryptStr, ce.secretKey)
}

func encrypt(value string, secretKey []byte) (string, error) {
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

func signature(body string, expire time.Time, hashKey []byte) string {
	msg := []byte(fmt.Sprintf("%s;%d;", body, expire.Unix()))
	h := hmac.New(sha256.New, hashKey)
	h.Write(msg)
	return base64.StdEncoding.EncodeToString(h.Sum(msg))
}

func verify(sgn []byte, hashKey []byte) (body string, err error) {
	msgLen := len(sgn) - 32
	if msgLen <= 0 {
		err = errors.New("signature too short")
		return
	}
	h := hmac.New(sha256.New, hashKey)
	h.Write(sgn[:msgLen])
	if !hmac.Equal(h.Sum(nil), sgn[msgLen:]) {
		err = errors.New("signature does not match")
		return
	}

	parts := bytes.SplitN(sgn, []byte(";"), 3)
	log.Printf("%s, %s",string(parts[0]), string(parts[1]))
	if len(parts) < 3 {
		err = errors.New("incorrect signature format")
		return
	}
	expUnix, err := strconv.ParseInt(string(parts[1]), 10, 64)
	if err != nil {
		return
	}
	expire := time.Unix(expUnix, 0)
	if !expire.IsZero() && time.Now().Compare(expire) > 0 {
		err = errors.New("cookie expired")
		return
	}
	return string(parts[0]), nil
}

func decrypt(cryptStr string, secretKey []byte) (decrypted string, err error) {
	cryptByte, err := base64.StdEncoding.DecodeString(cryptStr)
	if err != nil {
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
