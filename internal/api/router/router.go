package router

import (
	"net/http"
	"school-management/internal/api/handlers"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.RootHandler)

	// Teachers
	mux.HandleFunc("GET /teachers", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers", handlers.AddTeachersHandler)
	mux.HandleFunc("PATCH /teachers", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers", handlers.DeleteTeachersHandler)

	mux.HandleFunc("GET /teachers/{id}", handlers.GetOneTeacherHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeachersHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchOneTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)

	// Students
	mux.HandleFunc("GET /students", handlers.GetStudentsHandler)
	mux.HandleFunc("POST /students", handlers.AddStudentsHandler)
	mux.HandleFunc("PATCH /students", handlers.PatchStudentsHandler)
	mux.HandleFunc("DELETE /students", handlers.DeleteStudentsHandler)

	mux.HandleFunc("GET /students/{id}", handlers.GetOneStudentHandler)
	mux.HandleFunc("PUT /students/{id}", handlers.UpdateStudentsHandler)
	mux.HandleFunc("PATCH /students/{id}", handlers.PatchOneStudentHandler)
	mux.HandleFunc("DELETE /students/{id}", handlers.DeleteStudentHandler)

	mux.HandleFunc("GET /execs", handlers.ExecsHandler)

	return mux
}
