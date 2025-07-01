package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/demurk/tinyurl/internal/config"
	"github.com/demurk/tinyurl/internal/storage"
	"github.com/demurk/tinyurl/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var shortenTestCases = []struct {
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

func TestMain(m *testing.M) {
	config.Parse()
	os.Truncate(*config.FileStoragePath, 0)
	storage.Initialize()
	os.Exit(m.Run())
}

func TestShortage(t *testing.T) {
	for _, tc := range shortenTestCases {
		t.Run(tc.name, func(t *testing.T) {
			postRequest := httptest.NewRequest(http.MethodPost, *config.OriginURL, strings.NewReader(tc.fullURL))
			w := httptest.NewRecorder()
			postHandler := http.HandlerFunc(saveTextURLHandler)
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
			getHandler := http.HandlerFunc(getFullURLHandler)
			getHandler(ww, idRequest)
			idResult := ww.Result()
			defer idResult.Body.Close()

			assert.Equal(t, tc.fullURL, idResult.Header.Get("Location"))
			assert.Equal(t, http.StatusTemporaryRedirect, idResult.StatusCode)
		})
	}
}

func TestShortageJSON(t *testing.T) {
	for _, tc := range shortenTestCases {
		t.Run(tc.name, func(t *testing.T) {
			responseData := types.JSONPostRequestData{URL: tc.fullURL}
			jsonBody, _ := json.Marshal(responseData)

			postRequest := httptest.NewRequest(http.MethodPost, *config.OriginURL, bytes.NewReader(jsonBody))
			w := httptest.NewRecorder()
			postHandler := http.HandlerFunc(saveJSONURLHandler)
			postHandler(w, postRequest)
			result := w.Result()

			assert.Equal(t, http.StatusCreated, result.StatusCode)

			bodyBytes, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			var r types.JSONPostResponseData
			err = json.Unmarshal(bodyBytes, &r)
			require.NoError(t, err)

			assert.Equal(t, *config.ResultURL+"/"+tc.shortURL, r.Result)

			idRequest := httptest.NewRequest(http.MethodGet, *config.OriginURL, nil)
			idRequest.SetPathValue("id", tc.shortURL)
			ww := httptest.NewRecorder()
			getHandler := http.HandlerFunc(getFullURLHandler)
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
		for _, URL := range []string{"/", "/api/shorten"} {
			t.Run(tc.method, func(t *testing.T) {
				request := httptest.NewRequest(tc.method, URL, strings.NewReader("https://github.com/demurk/tinyurl"))
				w := httptest.NewRecorder()

				postHandler := http.HandlerFunc(saveTextURLHandler)
				postHandler(w, request)

				assert.Equal(t, tc.expectedCode, w.Code, "Invalid status code")
			})
		}
	}
}

func TestInvalidURL(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not an url"))
	w := httptest.NewRecorder()

	postHandler := http.HandlerFunc(saveTextURLHandler)
	postHandler(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid URL")
}

func TestGetHandler(t *testing.T) {
	t.Run(http.MethodGet, func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		getHandler := http.HandlerFunc(getFullURLHandler)
		getHandler(w, request)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
