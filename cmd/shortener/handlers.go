package main

import (
	"io"
	"net/http"

	"github.com/demurk/tinyurl/internal/config"
	"github.com/demurk/tinyurl/internal/db"
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

	fullURL := string(fullURLBytes)
	if !IsValidURL(fullURL) {
		http.Error(res, "Invalid URL", http.StatusBadRequest)
		return
	}

	shortURLId := setFullURL(string(fullURLBytes))
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(*config.ResultURL + "/" + shortURLId))
}

func getPage(res http.ResponseWriter, req *http.Request) {
	shortURL := req.PathValue("id")
	fullURL, err := getFullURL(shortURL)
	if err != nil {
		http.Error(res, "Url doesnt exists", http.StatusNotFound)
		return
	}
	http.Redirect(res, req, fullURL, http.StatusTemporaryRedirect)
}

func pingPage(res http.ResponseWriter, req *http.Request) {
	err := db.GetConnection().Ping()

	if err != nil {
		http.Error(res, "", http.StatusInternalServerError)
		return
	}
}
