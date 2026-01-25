package middlewares

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func Compression(next http.Handler) http.Handler {
	log.Println("Compression Middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(("Compression Middleware being returned"))
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Set the response header
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Wrap the ResponseWriter with gzipWriter
		wrappedWriter := &gzipResponseWriter{
			ResponseWriter: w,
			Writer:         gz,
		}

		next.ServeHTTP(wrappedWriter, r)
		fmt.Println(("Sent response from compression middleware"))
	})
}

// gzipResponseWriter wraps http.ResponseWriter to write gzipeed responses
type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}
