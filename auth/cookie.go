package auth

import (
	"fmt"
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

func NewCookieEncrypter() (*CookieEncrypter, error) {
	if keyLen := len(secretKey); keyLen == 0 {
		return nil, errors.New("secret key is not set")
	}
	return &CookieEncrypter{secretKey: secretKey}, nil
}

func (ce CookieEncrypter) Encode(c http.Cookie) (out http.Cookie, err error) {
	out = c
	// encrypt value
	enc, err := encrypt(c.Value, ce.secretKey)
	if err != nil {
		err = fmt.Errorf("could not encrypt: %w", err)
		return
	} 
	
	encoded := signature(enc, c.Expires, ce.secretKey)
	out.Value = encoded

	err = out.Valid()
	if err != nil {
		err = fmt.Errorf("cookie not valid: %w", err)
		return
	}
	return out, nil
}

func (ce CookieEncrypter) Decode(c http.Cookie) (value string, err error) {
	b, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		err = fmt.Errorf("could not decode cookie value: %w", err)
		return
	}

	cryptStr, err := verify(b, ce.secretKey)
	if err != nil {
		err = fmt.Errorf("failed to verify cookie: %w", err)
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
	if len(parts) < 3 {
		err = errors.New("incorrect signature format")
		return
	}
	expUnix, err := strconv.ParseInt(string(parts[1]), 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid expire part: %w", err)
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
