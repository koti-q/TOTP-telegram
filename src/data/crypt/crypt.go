package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func EncryptSecret(secret string, key string) (string, error) {
	s := []byte(secret)
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, s, nil)
	// Encode to base64 to ensure a valid UTF-8 string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptSecret(encodedSecret string, key string) (string, error) {
	k := []byte(key)
	// Decode secret from base64
	s, err := base64.StdEncoding.DecodeString(encodedSecret)
	if err != nil {
		return "", err
	}
	c, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(s) < nonceSize {
		return "", err
	}
	nonce, s := s[:nonceSize], s[nonceSize:]
	secret, err := gcm.Open(nil, nonce, s, nil)
	if err != nil {
		return "", err
	}
	return string(secret), nil
}
