package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// method based routing
	mux.HandleFunc("POST /items/create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "item created")
	})
	mux.HandleFunc("DELETE /items/create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "item deleted")
	})

	//Wildcard in pattern - path parameter
	mux.HandleFunc("/teachers/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Teacher id: %s", r.PathValue("id"))
	})

	//Wildcard with "..." pattern
	mux.HandleFunc("/files/{path...}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Path: %s", r.PathValue("path"))
	})

	mux.HandleFunc("/path1/{param1}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "param1: %s", r.PathValue("param1"))
	})

	// mux.HandleFunc("/{param2}/path2", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "param2: %s", r.PathValue("param2"))
	// })

	mux.HandleFunc("/path1/path2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "param1: %s", r.PathValue("param1"))
	})

	http.ListenAndServe(":8080", mux)
}
