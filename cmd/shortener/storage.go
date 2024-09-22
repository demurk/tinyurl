package main

import "errors"

type Repository map[string]string

var urlsStorage = make(Repository)

func getFullURL(shortURL string) (string, error) {
	fullURL, exists := urlsStorage[shortURL]
	if !exists {
		return "", errors.New("url doesnt exists")
	}
	return fullURL, nil
}

func setFullURL(fullURL string) string {
	shortURL := makeShortURL(fullURL)
	urlsStorage[shortURL] = fullURL
	return shortURL
}
