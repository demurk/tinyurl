package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/demurk/tinyurl/cmd/shortener/config"
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

	shortURLId := setFullURL(r.URL)
	responseData := PostResponseData{Result: *config.ResultURL + "/" + shortURLId}
	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		http.Error(res, "Error marshaling response JSON", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(jsonResponse)
}
