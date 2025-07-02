package urls_storage

import (
	"database/sql"
	"errors"

	"github.com/demurk/tinyurl/internal/db"
)

type User struct {
	ID    int
	Name  string
	Email string
}

func dbGetFullURL(shortURL string) (string, error) {
	conn := db.GetConnection()

	var fullURL string
	err := conn.QueryRow("SELECT full_url FROM urls WHERE short_url = $1", shortURL).Scan(&fullURL)
	if err == sql.ErrNoRows {
		return "", errors.New("URL doesnt exists")
	}
	return fullURL, nil
}

func dbSetFullURL(fullURL string) (string, error) {
	conn := db.GetConnection()

	shortURL := makeShortURL(fullURL)
	_, err := conn.Exec(`
		INSERT INTO urls (short_url, full_url) 
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING`,
		shortURL, fullURL,
	)

	if err != nil {
		return "", err
	}
	urlsStorage.Set(shortURL, fullURL)
	return shortURL, nil
}
