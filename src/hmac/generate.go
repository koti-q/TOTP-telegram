package totp_generator

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"time"
)

func Generate(secret string, timestamp int64) (string, error) {
	hmacHash := hmac.New(sha1.New, []byte(secret))

	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(timestamp))

	hmacHash.Write(timeBytes)
	hash := hmacHash.Sum(nil)

	offset := hash[len(hash)-1] & 0xf

	code := binary.BigEndian.Uint32(hash[offset : offset+4])
	code = code & 0x7fffffff // Remove most significant bit

	otp := code % 1000000

	return fmt.Sprintf("%06d", otp), nil
}

// Helper function to get current TOTP
func GenerateCurrentTOTP(secret string) (string, error) {
	timestamp := time.Now().Unix() / 30
	return Generate(secret, timestamp)
}
