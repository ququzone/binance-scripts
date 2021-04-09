package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type HmacSigner struct {
	secret []byte
}

func NewHmacSigner(secret string) *HmacSigner {
	return &HmacSigner{secret: []byte(secret)}
}

func (s HmacSigner) Sign(data []byte) string {
	h := hmac.New(sha256.New, s.secret)

	h.Write(data)

	return hex.EncodeToString(h.Sum(nil))
}
