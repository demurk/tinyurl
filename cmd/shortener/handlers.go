package main

import (
	"io"
	"net/http"

	"github.com/demurk/tinyurl/cmd/shortener/config"
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
	shortURLId := setFullURL(string(fullURLBytes))
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(*config.ResultURL + "/" + shortURLId))
}

func idPage(res http.ResponseWriter, req *http.Request) {
	shortURL := req.PathValue("id")
	fullURL, err := getFullURL(shortURL)
	if err != nil {
		http.Error(res, "Url doesnt exists", http.StatusNotFound)
		return
	}
	http.Redirect(res, req, fullURL, http.StatusTemporaryRedirect)
}
