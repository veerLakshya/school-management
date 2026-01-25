package middlewares

import (
	"log"
	"net/http"
	"slices"
)

// api is hosted at www.myapi.com
// frontend is hosted at www.myfrontend.com

// Allowed Origins
var allowedOrigins = []string{
	"https://my-origin-url.com",
	"https://www.myfrontend.com",
	"https://localhost:3000",
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		log.Printf("CORS middleware - Path: %s, Method: %s, Origin: %s\n", r.URL.Path, r.Method, origin)

		if isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			http.Error(w, "Cors Error", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Headers", "Control-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// to avoid preflight calls from triggering ahead
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isOriginAllowed(origin string) bool {
	// for _, allowedOrigin := range allowedOrigins {
	// 	if origin == allowedOrigin {
	// 		return true
	// 	}
	// }

	if slices.Contains(allowedOrigins, origin) {
		return true
	} else {
		return false
	}
}
