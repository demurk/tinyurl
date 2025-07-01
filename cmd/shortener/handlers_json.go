package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/demurk/tinyurl/internal/config"
	"github.com/demurk/tinyurl/internal/storage"
)

type PostRequestData struct {
	URL string `json:"url"`
}
type PostResponseData struct {
	Result string `json:"result"`
}

func postPageJSON(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	reqBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Couldn't read request body", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	var r PostRequestData
	err = json.Unmarshal(reqBytes, &r)
	if err != nil {
		http.Error(res, "Couldn't parse request json", http.StatusInternalServerError)
		return
	}

	if !IsValidURL(r.URL) {
		http.Error(res, "Invalid URL", http.StatusBadRequest)
		return
	}

	urlStorage := storage.Get()
	var shortURL string
	shortURL, err = urlStorage.Set(r.URL)
	if err != nil {
		http.Error(res, "Couldn't store url, try again", http.StatusInternalServerError)
		return
	}

	responseData := PostResponseData{Result: *config.ResultURL + "/" + shortURL}
	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		http.Error(res, "Error marshaling response JSON", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(jsonResponse)
}
