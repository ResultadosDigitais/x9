package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func SHA256(payload string) string {
	h := sha256.New()
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

func HMAC(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func ValidateHMAC(current, expected, prefix, secret string) bool {
	current = fmt.Sprintf("%s=%s", prefix, HMAC(current, secret))
	return current == expected
}
