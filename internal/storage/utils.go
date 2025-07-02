package storage

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/demurk/tinyurl/internal/config"
)

const shortURLLettersLimit = 8

func makeShortURL(fullURL string) string {
	hasher := sha256.New()
	hasher.Write([]byte(fullURL))
	hash := hasher.Sum(nil)
	return base64.URLEncoding.EncodeToString(hash)[:shortURLLettersLimit]
}

func ShortURLWithHost(shortURL string) string {
	return *config.ResultURL + "/" + shortURL
}
