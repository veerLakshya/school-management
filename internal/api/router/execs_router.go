package router

import (
	"net/http"
	"school-management/internal/api/handlers"
)

func execsRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /execs", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs", handlers.AddExecsHandler)
	mux.HandleFunc("PATCH /execs", handlers.PatchExecsHandler)

	mux.HandleFunc("GET /execs/{id}", handlers.GetOneExecHandler)
	mux.HandleFunc("PATCH /execs/{id}", handlers.PatchOneExecHandler)
	// mux.HandleFunc("POST /execs/{id}/updatepassword", handlers.AddExecsHandler)
	mux.HandleFunc("DELETE /execs/{id}", handlers.DeleteOneExecHandler)

	// mux.HandleFunc("POST /execs/login", handlers.ExecsLoginHandler)
	// mux.HandleFunc("POST /execs/logout", handlers.ExecsLogoutHandler)
	// mux.HandleFunc("POST /execs/forgotpassword", handlers.ExecsForgotPasswordHandler)
	// mux.HandleFunc("POST /execs/resetpassword/reset/{resetcode}", handlers.ExecsResetPasswordHandler)

	return mux
}
