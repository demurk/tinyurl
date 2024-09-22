package main

import (
	"crypto/sha256"
	"encoding/base64"
)

const shortUrlLettersLimit = 8

func makeShortURL(fullURL string) string {
	hasher := sha256.New()
	hasher.Write([]byte(fullURL))
	hash := hasher.Sum(nil)
	return base64.URLEncoding.EncodeToString(hash)[:shortUrlLettersLimit]
}
