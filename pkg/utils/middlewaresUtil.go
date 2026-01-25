package utils

import "net/http"

// Middleware is a function that properly wraps an http.Handler with additional functionality (for chaining without cluttering)
type Middleware func(http.Handler) http.Handler

func ApplyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
