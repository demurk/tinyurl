package compression

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

func GZIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			uncompressedBody, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decode body", http.StatusInternalServerError)
				return
			}
			uncompressedBody.Close()
			r.Body = uncompressedBody
		}

		acceptEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()

			next.ServeHTTP(&gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gz,
			}, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
