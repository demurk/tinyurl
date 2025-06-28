package middlewareLogger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func ConnectZapLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			next.ServeHTTP(rw, r)

			logger.Info("Request info",
				zap.String("URI", r.RequestURI),
				zap.String("method", r.Method),
				zap.Duration("duration_seconds", time.Since(start)),
			)
			logger.Info("Response info",
				zap.Int("status", rw.Status()),
				zap.Int("response_size_bytes", rw.BytesWritten()),
			)
		})
	}
}
