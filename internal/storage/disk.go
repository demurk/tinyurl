package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/demurk/tinyurl/internal/config"
	"github.com/demurk/tinyurl/internal/types"
)

var lastUUID = 0

type urlData struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func RestoreURLsFromFile() {
	file, err := os.OpenFile(*config.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var url urlData
		err := json.Unmarshal([]byte(line), &url)
		if err != nil {
			fmt.Printf("Error parsing string: %v\n", err)
			continue
		}

		mSetShortFullURL(url.ShortURL, url.OriginalURL)
		lastUUID = url.UUID
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	lastUUID += 1
}

func dSetFullURL(fullURL string) (string, error) {
	shortURL := makeShortURL(fullURL)
	jsonData, err := json.Marshal(
		urlData{UUID: lastUUID, ShortURL: shortURL, OriginalURL: fullURL},
	)
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(*config.FileStoragePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	file.Write(jsonData)
	file.WriteString("\n")

	mSetShortFullURL(shortURL, fullURL)

	lastUUID += 1
	return shortURL, nil
}

func dSetFullURLBatch(urlSlice []types.BatchJsonPostRequestData) ([]types.BatchJsonPostResponseData, error) {
	var returnValues []types.BatchJsonPostResponseData
	file, err := os.OpenFile(*config.FileStoragePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	for _, url := range urlSlice {
		shortURL := makeShortURL(url.OriginalURL)
		jsonData, err := json.Marshal(
			urlData{UUID: lastUUID, ShortURL: shortURL, OriginalURL: url.OriginalURL},
		)
		if err != nil {
			return nil, err
		}

		file.Write(jsonData)
		file.WriteString("\n")

		mSetShortFullURL(shortURL, url.OriginalURL)
		lastUUID += 1

		returnValues = append(returnValues, types.BatchJsonPostResponseData{
			CorrelationID: url.CorrelationID,
			ShortURL:      ShortURLWithHost(shortURL),
		})
	}

	return returnValues, nil
}
