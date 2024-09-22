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
	fullUrlBytes, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	shortUrl, err := setFullUrl(string(fullUrlBytes))
	if err != nil {
		http.Error(res, "Couldn't create short url, try again", http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortUrl))
}

func idPage(res http.ResponseWriter, req *http.Request) {
	shortUrl := req.PathValue("id")
	fullUrl, err := getFullUrl(shortUrl)
	if err != nil {
		http.Error(res, "Url doesnt exists", http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", "text/plain")
	http.Redirect(res, req, fullUrl, http.StatusTemporaryRedirect)
}
