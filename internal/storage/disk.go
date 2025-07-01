package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/demurk/tinyurl/internal/config"
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

	lastUUID += 1
	return shortURL, nil
}
