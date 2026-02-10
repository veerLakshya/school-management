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

// func OldTeachersHandler(w http.ResponseWriter, r *http.Request) {
// switch r.Method {
// case http.MethodGet:
// 	getTeachersHandler(w, r)
// case http.MethodPost:
// 	postTeachersHandler(w, r)
// case http.MethodPut:
// 	updateTeachersHandler(w, r)
// case http.MethodPatch:
// 	patchTeacherHandler(w, r)
// case http.MethodDelete:
// 	deleteTeacherHandler(w, r)
// }
// }

func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getTeacherHandler:", r.URL)

	idStr := r.PathValue("id")

	//Handle Path Parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing id: %v", err), http.StatusBadRequest)
		return
	}

	teacher, err := sqlconnect.GetTeacherByIdDbHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getTeachersHandler:", r.URL)

	var teachers []models.Teacher
	teachers, err := sqlconnect.GetTeachersDbHandler(teachers, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teachers),
		Data:   teachers,
	}
	w.Header().Set("X-Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AddTeachersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("postTeachersHandler:", r.URL, r.Body)

	reqBody, err := io.ReadAll(r.Body) // store req body as it gets wiped out on reading once only
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var newTeachers []models.Teacher
	var rawTeachers []map[string]interface{}

	err = json.Unmarshal(reqBody, &rawTeachers)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	fields := GetFieldNames(models.Teacher{})

	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	for _, teacher := range rawTeachers {
		for key := range teacher {
			_, ok := allowedFields[key]
			if !ok {
				http.Error(w, "Unacceptable field found in request. Only use allowed fields.", http.StatusBadRequest)
				return
			}
		}
	}

	err = json.Unmarshal(reqBody, &newTeachers)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	for _, teacher := range newTeachers {
		// if teacher.FirstName == "" || teacher.Email == "" || teacher.Class == "" || teacher.LastName == "" || teacher.Subject == "" {
		// 	http.Error(w, "All fields are required", http.StatusBadRequest)
		// 	return
		// }

		err = CheckBlankFields(teacher)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var addedTeachers []models.Teacher

	addedTeachers, err = sqlconnect.AddTeachersDbHandler(addedTeachers, newTeachers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}

	json.NewEncoder(w).Encode(response)

}

// PUT /teachers/{id}  - PUT replaces all fields even if left empty
func UpdateTeachersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("update teachers handler:", r.URL)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher id", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher

	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		http.Error(w, "Invalid Request payload", http.StatusBadRequest)
		return
	}

	updatedTeacherFromDB, err := sqlconnect.UpdateTeacherByIdDbHandle(id, updatedTeacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacherFromDB)
}

// PATCH /teachers/{id} - PATCH only updated the given fields
func PatchOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("patch teachers handler:", r.URL)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher id", http.StatusBadRequest)
		return
	}

	var updates map[string]string
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedTeacher, err := sqlconnect.PatchTeacherByIdDbHandler(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
}

// patch multiple teachers
func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchTeachersDbHandler(updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /teachers/{id}
func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("delete teacher handler:", r.URL)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher id", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteTeacherByIdDbHandler(id)

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
		Status: "Teacher successfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}

// DELETE /teachers - multiple teachers
func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("delete teachers handler:", r.URL)

	var ids []int

	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "Invalid teacher ids", http.StatusBadRequest)
		return
	}

	deletedIds, err := sqlconnect.DeleteTeachersDbHandler(ids)
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
		Status:     "Teachers successfully deleted",
		DeletedIds: deletedIds,
	}
	json.NewEncoder(w).Encode(response)

}

func GetStudentsByTeacherId(w http.ResponseWriter, r *http.Request) {
	teacherIdFromPath := r.PathValue("id")

	teacherId, err := strconv.Atoi(teacherIdFromPath)
	if err != nil {
		http.Error(w, "Invalid teacher Id", http.StatusBadRequest)
		return
	}

	var students []models.Student

	students, err = sqlconnect.GetStudentsByTeacherIdDbHandler(teacherId, students)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Data:   students,
		Count:  len(students),
		Status: "succuess",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetStudentCountById(w http.ResponseWriter, r *http.Request) {
	teacherIdFromPath := r.PathValue("id")

	teacherId, err := strconv.Atoi(teacherIdFromPath)
	if err != nil {
		http.Error(w, "Invalid teacher Id", http.StatusBadRequest)
		return
	}

	studentCount, err := sqlconnect.GetStudentCountByTeacherIdDbHandler(teacherId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}{
		Count:  studentCount,
		Status: "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
