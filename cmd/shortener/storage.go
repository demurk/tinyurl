package main

import "errors"

type Repository map[string]string

var urlsStorage = make(Repository)

func getFullUrl(shortUrl string) (string, error) {
	fullUrl, exists := urlsStorage[shortUrl]
	if exists {
		delete(urlsStorage, shortUrl)
	} else {
		return "", errors.New("url doesnt exists")
	}
	return fullUrl, nil
}

const urlShortAttempts = 10

func setFullUrl(fullUrl string) (string, error) {
	shortUrl := ""
	for i := range urlShortAttempts {
		shortUrl = makeRandomString()
		_, exists := urlsStorage[shortUrl]
		if !exists {
			break
		}
		if i == urlShortAttempts-1 {
			return "", errors.New("link shortening error")
		}
	}
	urlsStorage[shortUrl] = fullUrl
	return shortUrl, nil
}
