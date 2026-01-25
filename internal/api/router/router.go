package router

import (
	"net/http"
	"school-management/internal/api/handlers"
)

func NewRouter() *http.ServeMux {

	handlers.Init()

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.RootHandler)

	mux.HandleFunc("/teachers/", handlers.TeachersHandler)

	mux.HandleFunc("/students/", handlers.StudentsHandler)

	mux.HandleFunc("/execs", handlers.ExecsHandler)

	return mux
}
