package main

import (
	"io"
	"net/http"

	"github.com/demurk/tinyurl/internal/config"
	"github.com/demurk/tinyurl/internal/storage"
)

func postPage(res http.ResponseWriter, req *http.Request) {
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

	urlStorage := storage.Get()
	var shortURL string
	shortURL, err = urlStorage.Set(string(fullURLBytes))
	if err != nil {
		http.Error(res, "Couldn't store url, try again", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(*config.ResultURL + "/" + shortURL))
}

func getPage(res http.ResponseWriter, req *http.Request) {
	shortURL := req.PathValue("id")
	urlStorage := storage.Get()
	fullURL, err := urlStorage.Get(shortURL)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(res, req, fullURL, http.StatusTemporaryRedirect)
}
