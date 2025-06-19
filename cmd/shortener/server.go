package main

import (
	"net/http"

	"github.com/demurk/tinyurl/cmd/shortener/config"
	"github.com/demurk/tinyurl/internal/app/middleware/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	config.Parse()

	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync()

	r := chi.NewRouter()

	r.Use(logger.ZapLogger(zapLogger))
	r.MethodNotAllowed(func(res http.ResponseWriter, r *http.Request) {
		http.Error(res, "Invalid request method", http.StatusMethodNotAllowed)
	})

	r.Post("/", postPage)
	r.Get("/{id}", getPage)

	r.Post("/api/shorten", postPageJSON)

	err = http.ListenAndServe(*config.OriginURL, r)
	if err != nil {
		panic(err)
	}
}
