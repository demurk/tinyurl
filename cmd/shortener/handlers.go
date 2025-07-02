package main

import (
	"errors"
	"io"
	"net/http"

	"github.com/demurk/tinyurl/internal/db"
	"github.com/demurk/tinyurl/internal/storage"
)

func saveTextURLHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	fullURLBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Couldn't read request body", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	fullURL := string(fullURLBytes)
	if !IsValidURL(fullURL) {
		http.Error(res, "Invalid URL", http.StatusBadRequest)
		return
	}

	urlStorage := storage.Get()
	var shortURL string
	shortURL, err = urlStorage.Set(fullURL)
	if err != nil {
		if errors.Is(err, storage.ErrURLAlreadyExists) {
			res.WriteHeader(http.StatusConflict)
		} else {
			http.Error(res, "Couldn't store url, try again", http.StatusInternalServerError)
			return
		}
	} else {
		res.WriteHeader(http.StatusCreated)
	}
	res.Header().Set("content-type", "text/plain")
	res.Write([]byte(storage.ShortURLWithHost(shortURL)))
}

func getFullURLHandler(res http.ResponseWriter, req *http.Request) {
	shortURL := req.PathValue("id")
	urlStorage := storage.Get()
	fullURL, err := urlStorage.Get(shortURL)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(res, req, fullURL, http.StatusTemporaryRedirect)
}

func pingDBPage(res http.ResponseWriter, req *http.Request) {
	err := db.GetConnection().Ping()

	if err != nil {
		http.Error(res, "", http.StatusInternalServerError)
		return
	}
}
