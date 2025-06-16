package main

import (
	"fmt"
	"net/http"

	"github.com/demurk/tinyurl/cmd/shortener/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.ParseFlags()
	r := chi.NewRouter()
	r.MethodNotAllowed(func(res http.ResponseWriter, r *http.Request) {
		http.Error(res, "Invalid request method", http.StatusMethodNotAllowed)
	})
	r.Use(middleware.Logger)
	r.Post("/", postPage)
	r.Get(fmt.Sprintf("/%s{id}", *config.BaseResultURL), idPage)

	err := http.ListenAndServe(*config.OriginURL, r)
	if err != nil {
		panic(err)
	}
}
