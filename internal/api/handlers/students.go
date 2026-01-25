package handlers

import (
	"log"
	"net/http"
)

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/students  handler", r.Method)

	w.Write([]byte("Hello Students Route"))
}
