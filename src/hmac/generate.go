package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"strings"
)

func generateTOTP(secretKey string, timestamp int64) uint32 {
	base32Decoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	secretKey = strings.ToUpper(strings.TrimSpace(secretKey))
	secretBytes, _ := base32Decoder.DecodeString(secretKey)

	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(timestamp)/30)

	hash := hmac.New(sha1.New, secretBytes)
	hash.Write(timeBytes)
	h := hash.Sum(nil)

	offset := h[len(h)-1] & 0x0F
	truncatedHash := binary.BigEndian.Uint32(h[offset:]) & 0x7FFFFFFF

	return truncatedHash % 1_000_000
}
