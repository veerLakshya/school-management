package handlers

import (
	"log"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ROOT Route", r.Method)
	// fmt.Fprintf(w, "Hello Root Route")
	w.Write([]byte("Welcome!"))

}
