package main

import (
	"crypto/sha256"
	"encoding/base64"
	"net/url"
)

const shortURLLettersLimit = 8

func makeShortURL(fullURL string) string {
	hasher := sha256.New()
	hasher.Write([]byte(fullURL))
	hash := hasher.Sum(nil)
	return base64.URLEncoding.EncodeToString(hash)[:shortURLLettersLimit]
}

func IsValidURL(urlString string) bool {
	u, err := url.Parse(urlString)
	return err == nil && u.Scheme != "" && u.Host != ""
}
