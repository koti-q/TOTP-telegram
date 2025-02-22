package data

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/bcrypt"
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
	return string(ciphertext), nil
}

func DecryptSecret(c_secret string, key string) (string, error) {
	k := []byte(key)
	s := []byte(c_secret)
	c, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, s := s[:nonceSize], s[nonceSize:]
	secret, err := gcm.Open(nil, nonce, s, nil)
	if err != nil {
		return "", err
	}
	return string(secret), nil
}

// Hashing user key that can decrypt the secret
func HashKey(key string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(key), 14)
	return string(bytes), err
}

func CompareKey(key, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(key))
	return err == nil
}
