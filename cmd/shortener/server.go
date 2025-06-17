package main

import (
	"net/http"

	"github.com/demurk/tinyurl/cmd/shortener/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.Parse()
	r := chi.NewRouter()
	r.MethodNotAllowed(func(res http.ResponseWriter, r *http.Request) {
		http.Error(res, "Invalid request method", http.StatusMethodNotAllowed)
	})
	r.Use(middleware.Logger)
	r.Post("/", postPage)
	r.Get("/{id}", getPage)

	err := http.ListenAndServe(*config.OriginURL, r)
	if err != nil {
		panic(err)
	}
}
