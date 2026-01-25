package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"school-management/internal/models"
	"strconv"
	"strings"
	"sync"
)

var (
	teachers = make(map[int]models.Teacher)
	mutex    = &sync.Mutex{}
	nextID   = 1
)

// Initialize some data
func Init() {
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "9A",
		Subject:   "Math",
	}
	nextID++
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "10A",
		Subject:   "Science",
	}
	nextID++
}

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		postTeachersHandler(w, r)
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE on Teachers Route"))
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH on Teachers Route"))
	case http.MethodPut:
		w.Write([]byte("Hello PUT on Teachers Route"))
	}

	// fmt.Fprintf(w, "Hello Root Route")
	// w.Write([]byte("Hello Teachers Route"))
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getTeachersHandler:", r.URL)

	teacherList := make([]models.Teacher, 0, len(teachers))

	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")

	// Handle Query Parameters
	if idStr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")
		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
				teacherList = append(teacherList, teacher)
			}
		}
		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}
		w.Header().Set("X-Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	//Handle Path Parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error parsing id:", err)
		http.Error(w, fmt.Sprintf("Error parsing id: %v", err), http.StatusBadRequest)
		return
	}

	teacher, exists := teachers[id]
	if !exists {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(teacher)

}

func postTeachersHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	log.Println("postTeachersHandler:", r.URL, r.Body)

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)

	if err != nil {
		http.Error(w, fmt.Sprintf(`Invalid Request Body: %v`, err), http.StatusBadRequest)
		return
	}
	fmt.Println(newTeachers)

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		newTeacher.ID = nextID
		teachers[nextID] = newTeacher
		addedTeachers[i] = newTeacher
		nextID++
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
