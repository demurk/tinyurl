package main

import (
	"io"
	"net/http"
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
	shortURLId, err := setFullURL(string(fullURLBytes))
	if err != nil {
		http.Error(res, "Couldn't create short url, try again", http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("http://localhost:8080/" + shortURLId))
}

func idPage(res http.ResponseWriter, req *http.Request) {
	shortURL := req.PathValue("id")
	fullURL, err := getFullURL(shortURL)
	if err != nil {
		http.Error(res, "Url doesnt exists", http.StatusNotFound)
		return
	}
	res.Header().Set("Location", "text/plain")
	http.Redirect(res, req, fullURL, http.StatusTemporaryRedirect)
}
