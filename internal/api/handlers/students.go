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

// func OldStudentsHandler(w http.ResponseWriter, r *http.Request) {
// switch r.Method {
// case http.MethodGet:
// 	getStudentsHandler(w, r)
// case http.MethodPost:
// 	postStudentsHandler(w, r)
// case http.MethodPut:
// 	updateStudentsHandler(w, r)
// case http.MethodPatch:
// 	patchStudentHandler(w, r)
// case http.MethodDelete:
// 	deleteStudentHandler(w, r)
// }
// }

func GetOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getStudentHandler:", r.URL)

	idStr := r.PathValue("id")

	//Handle Path Parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing id: %v", err), http.StatusBadRequest)
		return
	}

	Student, err := sqlconnect.GetStudentByIdDbHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(Student)
}

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getStudentsHandler:", r.URL)

	var Students []models.Student
	Students, err := sqlconnect.GetStudentsDbHandler(Students, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(Students),
		Data:   Students,
	}
	w.Header().Set("X-Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AddStudentsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("postStudentsHandler:", r.URL, r.Body)

	reqBody, err := io.ReadAll(r.Body) // store req body as it gets wiped out on reading once only
	fmt.Println(string(reqBody))
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var newStudents []models.Student
	var rawStudents []map[string]interface{}

	err = json.Unmarshal(reqBody, &rawStudents)
	if err != nil {
		http.Error(w, "Invalid Request Bodyas", http.StatusBadRequest)
		return
	}

	fields := GetFieldNames(models.Student{})

	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	for _, Student := range rawStudents {
		for key := range Student {
			_, ok := allowedFields[key]
			if !ok {
				http.Error(w, "Unacceptable field found in request. Only use allowed fields.", http.StatusBadRequest)
				return
			}
		}
	}

	err = json.Unmarshal(reqBody, &newStudents)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	for _, Student := range newStudents {
		// if Student.FirstName == "" || Student.Email == "" || Student.Class == "" || Student.LastName == "" || Student.Subject == "" {
		// 	http.Error(w, "All fields are required", http.StatusBadRequest)
		// 	return
		// }

		err = CheckBlankFields(Student)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var addedStudents []models.Student

	addedStudents, err = sqlconnect.AddStudentsDbHandler(addedStudents, newStudents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(addedStudents),
		Data:   addedStudents,
	}

	json.NewEncoder(w).Encode(response)

}

// PUT /Students/{id}  - PUT replaces all fields even if left empty
func UpdateStudentsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("update Students handler:", r.URL)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Student id", http.StatusBadRequest)
		return
	}

	var updatedStudent models.Student

	err = json.NewDecoder(r.Body).Decode(&updatedStudent)
	if err != nil {
		http.Error(w, "Invalid Request payload", http.StatusBadRequest)
		return
	}

	updatedStudentFromDB, err := sqlconnect.UpdateStudentByIdDbHandle(id, updatedStudent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStudentFromDB)
}

// PATCH /Students/{id} - PATCH only updated the given fields
func PatchOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("patch Students handler:", r.URL)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Student id", http.StatusBadRequest)
		return
	}

	var updates map[string]string
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedStudent, err := sqlconnect.PatchStudentByIdDbHandler(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStudent)
}

// patch multiple Students
func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchStudentsDbHandler(updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /Students/{id}
func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("delete Student handler:", r.URL)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Student id", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteStudentByIdDbHandler(id)

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
		Status: "Student successfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}

// DELETE /Students - multiple Students
func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("delete Students handler:", r.URL)

	var ids []int

	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "Invalid Student ids", http.StatusBadRequest)
		return
	}

	deletedIds, err := sqlconnect.DeleteStudentsDbHandler(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//response body
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status     string `json:"status"`
		DeletedIds []int  `json:"deleted_ids"`
	}{
		Status:     "Students successfully deleted",
		DeletedIds: deletedIds,
	}
	json.NewEncoder(w).Encode(response)

}
