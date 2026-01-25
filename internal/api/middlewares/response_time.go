package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func ResponseTime(next http.Handler) http.Handler {
	fmt.Println("Response Time Middleware ...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Response Time Middleware being returned ...")
		start := time.Now()

		// create a custom ResponseWriter to capture the status code
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(wrappedWriter, r)

		// calculate the duration and log it
		duration := time.Since(start)
		wrappedWriter.Header().Set("X-Response-Time", duration.String()) // needs to be done before handing over control to next
		log.Printf("Method: %s, URL: %s, Status: %d, Duration: %v\n", r.Method, r.URL, wrappedWriter.status, duration.String())
		fmt.Println("Response Time Middleware ends ...")
	})
}

// responseWriter
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
