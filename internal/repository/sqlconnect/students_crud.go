package sqlconnect

import (
	"database/sql"
	"net/http"
	"reflect"
	"school-management/internal/models"
	"school-management/pkg/utils"
	"strconv"
	"strings"
)

func GetStudentByIdDbHandler(id int) (models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Error connecting to DB")
	}
	defer db.Close()

	var student models.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)

	if err == sql.ErrNoRows {
		return models.Student{}, utils.ErrorHandler(err, "Student Not found")
	} else if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Error retrieving Student")
	}
	return student, nil
}

func GetStudentsDbHandler(students []models.Student, r *http.Request) ([]models.Student, error) {

	query := "SELECT id, first_name, last_name, email, class FROM students WHERE 1 = 1"
	var args []interface{}

	query, args = utils.AddFilters(r, query, args)

	// also handling - Students/?sortby=name:asc&sortby=class:desc
	query = utils.AddSorting(r, query)

	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error querying DB")
	}

	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Class, &student.Email)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error scaning row from db")
		}
		students = append(students, student)
	}
	return students, nil
}

func AddStudentsDbHandler(addedStudents []models.Student, newStudents []models.Student) ([]models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	// stmt, err := db.Prepare("INSERT INTO students (first_name, last_name, email, class) VALUES(?,?,?,?,?)")
	stmt, err := db.Prepare(utils.GenerateInsertQuery("Students", models.Student{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error updating Student")
	}
	defer stmt.Close()

	for _, newStudent := range newStudents {
		values := utils.GetStructValues(newStudent)
		// res, err := stmt.Exec(newStudent.FirstName, newStudent.LastName, newStudent.Email, newStudent.Class)
		res, err := stmt.Exec(values...)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				return nil, utils.ErrorHandler(err, "")
			}
			return nil, utils.ErrorHandler(err, "Error updating Student")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error updating Student")
		}
		newStudent.ID = int(lastID)
		addedStudents = append(addedStudents, newStudent)
	}
	return addedStudents, nil
}

func UpdateStudentByIdDbHandle(id int, updatedStudent models.Student) (models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	var existingStudent models.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(&existingStudent.ID, &existingStudent.FirstName, &existingStudent.LastName, &existingStudent.Email, &existingStudent.Class)
	if err == sql.ErrNoRows {
		return models.Student{}, utils.ErrorHandler(err, "Student not found")
	} else if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Error updating Student")
	}

	updatedStudent.ID = existingStudent.ID
	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?", updatedStudent.FirstName, updatedStudent.LastName, updatedStudent.Email, updatedStudent.Class, updatedStudent.ID)
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Error updating Student")
	}
	return updatedStudent, nil
}

func DeleteStudentByIdDbHandler(id int) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "Error deleting Student")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "Error deleting Student")
	}

	if rowsAffected == 0 {
		return utils.ErrorHandler(err, "Student not found")
	}
	return nil
}

func PatchStudentsDbHandler(updates []map[string]interface{}) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	// transactions are used for commands which should either execute all or all fail
	tx, err := db.Begin()
	if err != nil {
		return utils.ErrorHandler(err, "Error starting transaction")
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			return utils.ErrorHandler(err, "Error retrieving id")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Invalid Student Id")
		}

		var studentFromDb models.Student
		err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(&studentFromDb.ID, &studentFromDb.FirstName, &studentFromDb.LastName, &studentFromDb.Email, &studentFromDb.Class)

		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				return utils.ErrorHandler(err, "Student not found")
			}
			return utils.ErrorHandler(err, "Error updating Students")
		}

		studentVal := reflect.ValueOf(&studentFromDb).Elem()
		studentType := studentVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < studentVal.NumField(); i++ {
				field := studentType.Field(i)
				// if field.Tag.Get("json") == k+",omitempty" {
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := studentVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							return utils.ErrorHandler(err, "Error getting filed value")
						}
					}
					break
				}
			}
		}
		_, err = tx.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?", studentFromDb.FirstName, studentFromDb.LastName, studentFromDb.Email, studentFromDb.Class, studentFromDb.ID)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Error updating Students")
		}
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return utils.ErrorHandler(err, "Error commiting transaction")
	}
	return nil
}

func PatchStudentByIdDbHandler(id int, updates map[string]string) (models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	var existingStudent models.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students where id = ?", id).Scan(&existingStudent.ID, &existingStudent.FirstName, &existingStudent.LastName, &existingStudent.Email, &existingStudent.Class)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Student{}, utils.ErrorHandler(err, "Student not found")
		}
		return models.Student{}, utils.ErrorHandler(err, "Error updating Student")
	}

	studentVal := reflect.ValueOf(&existingStudent).Elem()
	studentType := studentVal.Type()

	for k, v := range updates {
		for i := 0; i < studentVal.NumField(); i++ {
			field := studentType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if studentVal.Field(i).CanSet() {
					fieldVal := studentVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(studentVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?", existingStudent.FirstName, existingStudent.LastName, existingStudent.Email, existingStudent.Class, existingStudent.ID)
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Error updating Student")
	}
	return existingStudent, nil
}

func DeleteStudentsDbHandler(ids []int) ([]int, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error starting transaction")
	}

	stmt, err := tx.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "Error preparing query")
	}
	defer stmt.Close()

	deletedIds := []int{}

	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Error executing query")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Error getting affected rows")
		}

		if rowsAffected == 0 {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Student not found")
		} else if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
	}

	if len(deletedIds) != len(ids) {
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "Error deleting Students")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "Error commiting deleted changes")
	}
	return deletedIds, nil
}
