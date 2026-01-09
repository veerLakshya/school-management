package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/  handler", r.Method)

	// fmt.Fprintf(w, "Hello Root Route")
	w.Write([]byte("Hello Root Route"))

}

func teachersHandler(w http.ResponseWriter, r *http.Request) {

	// PATH PARAMS --> teachers/{id}
	fmt.Println("/teachers  handler", r.Method, r.URL.Path)
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	userId := strings.TrimSuffix(path, "/")
	fmt.Println("The id is: ", userId)

	// QUERY PARAMS --> teachers/?key=value&query=value2&sortby=email&sortorder=ASC
	fmt.Println("Query Params: ", r.URL.Query())
	queryParams := r.URL.Query()
	sortBy := queryParams.Get("sortby")
	sortOrder := queryParams.Get("sortorder")
	fmt.Println("Sort By: ", sortBy)
	fmt.Println("Sort Order: ", sortOrder)

	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET on Teachers Route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE on Teachers Route"))
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH on Teachers Route"))
	case http.MethodPost:
		w.Write([]byte("Hello POST on Teachers Route"))
	case http.MethodPut:
		w.Write([]byte("Hello PUT on Teachers Route"))
	}

	// fmt.Fprintf(w, "Hello Root Route")
	// w.Write([]byte("Hello Teachers Route"))
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/students  handler", r.Method)

	// fmt.Fprintf(w, "Hello Root Route")
	w.Write([]byte("Hello Students Route"))
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/execs  handler", r.Method)

	// fmt.Fprintf(w, "Hello Root Route")
	w.Write([]byte("Hello Execs Route"))
}

func main() {
	port := ":3000"

	http.HandleFunc("/", rootHandler)

	http.HandleFunc("/teachers/", teachersHandler)

	http.HandleFunc("/students/", studentsHandler)

	http.HandleFunc("/execs", execsHandler)

	fmt.Println("Server is running on port: ", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("Error starting the server: ", err)
	}
}
