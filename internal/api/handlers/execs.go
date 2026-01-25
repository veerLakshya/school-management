package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/execs handler", r.Method)

	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET on Execs Route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE on Execs Route"))
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH on Execs Route"))
	case http.MethodPost:
		fmt.Println("Query:", r.URL.Query())
		fmt.Println("Query:", r.URL.Query().Get("name"))

		//Parse form data(necessary for x-www-form-urlencoded)
		err := r.ParseForm()
		if err != nil {
			return
		}
		fmt.Println("Form from POST method:", r.Form)
		w.Write([]byte("Hello POST on Execs Route"))
	case http.MethodPut:
		w.Write([]byte("Hello PUT on Execs Route"))
	}
}
