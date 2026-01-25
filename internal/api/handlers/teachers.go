package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"school-management/internal/models"
	"school-management/internal/repository/sqlconnect"
	"strconv"
	"strings"
)

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

	db, err := sqlconnect.ConnectDB()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error connection to database: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")

	// Handle Query Parameters
	if idStr == "" {

		query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1 = 1"
		var args []interface{}

		// if firstName != "" {
		// 	query += " AND first_name = ?"
		// 	args = append(args, firstName)
		// }
		// if lastName != "" {
		// 	query += " AND last_name = ?"
		// 	args = append(args, lastName)
		// }

		query, args = addFilters(r, query, args)

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving data: %v", err), http.StatusInternalServerError)
			return
		}

		teacherList := make([]models.Teacher, 0)
		for rows.Next() {
			var teacher models.Teacher
			err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Class, &teacher.Email, &teacher.Subject)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
				return
			}
			teacherList = append(teacherList, teacher)
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

	var teacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Error finding teahcer: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func addFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"subject":    "subject",
		"class":      "class",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}
	return query, args
}

func postTeachersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("postTeachersHandler:", r.URL, r.Body)
	db, err := sqlconnect.ConnectDB()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error connection to database: %v", err), http.StatusInternalServerError)
		return
	}

	defer db.Close()

	var newTeachers []models.Teacher

	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, fmt.Sprintf(`Invalid Request Body: %v`, err), http.StatusBadRequest)
		return
	}
	fmt.Println(newTeachers)

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES(?,?,?,?,?)")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error preparing query: %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error inserting data into database: %v", err), http.StatusInternalServerError)
			return
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting last insert Id: %v", err), http.StatusInternalServerError)
			return
		}
		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher

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
