package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/demurk/tinyurl/cmd/shortener/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortage(t *testing.T) {
	config.Parse()

	testCases := []struct {
		name       string
		fullURL    string
		shortURL   string
		statusCode int
	}{
		{
			name:     "test shortage #1",
			fullURL:  "https://practicum.yandex.ru/learn/go-advanced",
			shortURL: "kpEZc4zh",
		},
		{
			name:     "test shortage #2",
			fullURL:  "https://go.dev/doc/tutorial/getting-started",
			shortURL: "MQ4QfaXg",
		},
		{
			name:     "test shortage #3",
			fullURL:  "https://github.com/demurk/tinyurl",
			shortURL: "NrBfyOZ7",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			postRequest := httptest.NewRequest(http.MethodPost, *config.OriginURL, strings.NewReader(tc.fullURL))
			w := httptest.NewRecorder()
			postHandler := http.HandlerFunc(postPage)
			postHandler(w, postRequest)
			result := w.Result()

			assert.Equal(t, http.StatusCreated, result.StatusCode)

			shortURLBytes, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			shortURLString := string(shortURLBytes)

			assert.Equal(t, *config.ResultURL+"/"+tc.shortURL, shortURLString)

			idRequest := httptest.NewRequest(http.MethodGet, *config.OriginURL, nil)
			idRequest.SetPathValue("id", tc.shortURL)
			ww := httptest.NewRecorder()
			getHandler := http.HandlerFunc(getPage)
			getHandler(ww, idRequest)
			idResult := ww.Result()
			defer idResult.Body.Close()

			assert.Equal(t, tc.fullURL, idResult.Header.Get("Location"))
			assert.Equal(t, http.StatusTemporaryRedirect, idResult.StatusCode)
		})
	}
}

func TestPostHandlerMethods(t *testing.T) {
	testCases := []struct {
		method       string
		expectedCode int
	}{
		{method: http.MethodGet, expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodPost, expectedCode: http.StatusCreated},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, "/", strings.NewReader("https://github.com/demurk/tinyurl"))
			w := httptest.NewRecorder()

			postHandler := http.HandlerFunc(postPage)
			postHandler(w, request)

			assert.Equal(t, tc.expectedCode, w.Code, "Invalid status code")
		})
	}
}

func TestInvalidURL(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not an url"))
	w := httptest.NewRecorder()

	postHandler := http.HandlerFunc(postPage)
	postHandler(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid URL")
}

func TestGetHandler(t *testing.T) {
	t.Run(http.MethodGet, func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		getHandler := http.HandlerFunc(getPage)
		getHandler(w, request)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
