package main

import (
	"net/http"

	"github.com/demurk/tinyurl/internal/app/middleware/compression"
	middlewarelogger "github.com/demurk/tinyurl/internal/app/middleware/logger"
	"github.com/demurk/tinyurl/internal/config"
	"github.com/demurk/tinyurl/internal/db"
	"github.com/demurk/tinyurl/internal/logger"
	"github.com/demurk/tinyurl/internal/urls_storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	logger.InitLogger()
	zapLogger := logger.GetLogger()
	defer zapLogger.Sync()

	config.Parse()

	db.Connect()
	conn := db.GetConnection()
	defer conn.Close()

	urls_storage.New()

	r := chi.NewRouter()

	r.Use(middlewarelogger.ConnectZapLogger(zapLogger))
	r.Use(compression.GZIP)
	r.MethodNotAllowed(func(res http.ResponseWriter, r *http.Request) {
		http.Error(res, "Invalid request method", http.StatusMethodNotAllowed)
	})

	r.Post("/", postPage)
	r.Get("/{id}", getPage)
	r.Get("/ping", pingPage)

	r.Post("/api/shorten", postPageJSON)

	err := http.ListenAndServe(*config.OriginURL, r)
	if err != nil {
		panic(err)
	}
}
