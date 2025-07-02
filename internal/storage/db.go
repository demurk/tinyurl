package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/demurk/tinyurl/internal/db"
	"github.com/demurk/tinyurl/internal/types"
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
		ON CONFLICT DO NOTHING;`,
		shortURL, fullURL,
	)

	if err != nil {
		return "", err
	}
	urlsStorage.Set(shortURL, fullURL)
	return shortURL, nil
}

func dbSetFullURLBatch(urlSlice []types.BatchJSONPostRequestData) ([]types.BatchJSONPostResponseData, error) {
	batchSize := 500
	conn := db.GetConnection()
	var returnValues []types.BatchJSONPostResponseData

	tx, _ := conn.Begin()
	var err error
	for i := 0; i < len(urlSlice); i += batchSize {
		endMarker := i + batchSize
		if endMarker > len(urlSlice) {
			endMarker = len(urlSlice)
		}
		batch := urlSlice[i:endMarker]

		query := "INSERT INTO urls (short_url, full_url) VALUES "

		values := []interface{}{}
		placeholders := []string{}

		for j := 0; j < len(batch); j += 1 {
			shortURL := makeShortURL(batch[j].OriginalURL)

			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", j*2+1, j*2+2))
			values = append(values, shortURL, batch[j].OriginalURL)

			returnValues = append(returnValues, types.BatchJSONPostResponseData{
				CorrelationID: batch[j].CorrelationID,
				ShortURL:      ShortURLWithHost(shortURL),
			})
		}

		query += strings.Join(placeholders, ",")
		query += " ON CONFLICT DO NOTHING;"

		_, err = tx.Exec(query, values...)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	err = tx.Commit()
	return returnValues, err
}

func CreateTables() {
	conn := db.GetConnection()
	conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.urls (
			id serial4 NOT NULL,
			full_url text NOT NULL,
			short_url varchar NOT NULL,
			CONSTRAINT urls_pk PRIMARY KEY (id),
			CONSTRAINT urls_un UNIQUE (full_url)
		);
	`)
}
