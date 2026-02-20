package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"school-management/internal/models"
	"school-management/internal/repository/sqlconnect"
	"strconv"
)

func GetOneExecHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getOneExecHandler:", r.URL)

	idStr := r.PathValue("id")

	//Handle Path Parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing id: %v", err), http.StatusBadRequest)
		return
	}

	exec, err := sqlconnect.GetExecByIdDbHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(exec)
}

func GetExecsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getExecsHandler:", r.URL)

	var execs []models.Exec
	execs, err := sqlconnect.GetExecsDbHandler(execs, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(execs),
		Data:   execs,
	}
	w.Header().Set("X-Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AddExecsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("postExecsHandler:", r.URL, r.Body)

	reqBody, err := io.ReadAll(r.Body) // store req body as it gets wiped out on reading once only
	fmt.Println(string(reqBody))
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var newExecs []models.Exec
	var rawExecs []map[string]interface{}

	err = json.Unmarshal(reqBody, &rawExecs)
	if err != nil {
		http.Error(w, "Invalid Request Bodyas", http.StatusBadRequest)
		return
	}

	fields := GetFieldNames(models.Exec{})

	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	for _, exec := range rawExecs {
		for key := range exec {
			_, ok := allowedFields[key]
			if !ok {
				http.Error(w, "Unacceptable field found in request. Only use allowed fields.", http.StatusBadRequest)
				return
			}
		}
	}

	err = json.Unmarshal(reqBody, &newExecs)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	for _, exec := range newExecs {
		err = CheckBlankFields(exec)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var addedExecs []models.Exec

	addedExecs, err = sqlconnect.AddExecsDbHandler(addedExecs, newExecs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(addedExecs),
		Data:   addedExecs,
	}

	json.NewEncoder(w).Encode(response)

}

// PATCH /execs/{id} - PATCH only updated the given fields
func PatchOneExecHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("patch Execss handler:", r.URL)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Exec id", http.StatusBadRequest)
		return
	}

	var updates map[string]string
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedExec, err := sqlconnect.PatchExecByIdDbHandler(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedExec)
}

// patch multiple Execss
func PatchExecsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Printf("error decoding: %s", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchExecsDbHandler(updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /execs/{id}
func DeleteOneExecHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("delete Exec handler:", r.URL)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Exec id", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteExecByIdDbHandler(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ---- alternate response ----
	// w.WriteHeader(http.StatusNoContent)

	//response body
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Exec successfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}
