package main

import (
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
	shortURL, err = urlStorage.Set(string(fullURLBytes))
	if err != nil {
		http.Error(res, "Couldn't store url, try again", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(storage.ShortURLWithHost(shortURL)))
}

func getFullURLHandler(res http.ResponseWriter, req *http.Request) {
	shortURL := req.PathValue("id")
	storage := storage.Get()
	fullURL, err := storage.Get(shortURL)
	if err != nil {
		http.Error(res, "Url doesnt exists", http.StatusNotFound)
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
