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

func GetTeacherByIdDbHandler(id int) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Error connecting to DB")
	}
	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "Teacher Not found")
	} else if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Error retrieving teacher")
	}
	return teacher, nil
}

func GetTeachersDbHandler(teachers []models.Teacher, r *http.Request) ([]models.Teacher, error) {

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1 = 1"
	var args []interface{}

	query, args = utils.AddFilters(r, query, args)

	// also handling - teachers/?sortby=name:asc&sortby=class:desc
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
		var teacher models.Teacher
		err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Class, &teacher.Email, &teacher.Subject)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error scaning row from db")
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func AddTeachersDbHandler(addedTeachers []models.Teacher, newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	// stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES(?,?,?,?,?)")
	stmt, err := db.Prepare(utils.GenerateInsertQuery("teachers", models.Teacher{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error updating teacher")
	}
	defer stmt.Close()

	for _, newTeacher := range newTeachers {
		values := utils.GetStructValues(newTeacher)
		// res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		res, err := stmt.Exec(values...)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				return nil, utils.ErrorHandler(err, "")
			}
			return nil, utils.ErrorHandler(err, "Error updating teacher")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error updating teacher")
		}
		newTeacher.ID = int(lastID)
		addedTeachers = append(addedTeachers, newTeacher)
	}
	return addedTeachers, nil
}

func UpdateTeacherByIdDbHandle(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "Teacher not found")
	} else if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Error updating teacher")
	}

	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Error updating teacher")
	}
	return updatedTeacher, nil
}

func DeleteTeacherByIdDbHandler(id int) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "Error deleting teacher")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "Error deleting teacher")
	}

	if rowsAffected == 0 {
		return utils.ErrorHandler(err, "Teacher not found")
	}
	return nil
}

func PatchTeachersDbHandler(updates []map[string]interface{}) error {
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
			return utils.ErrorHandler(err, "Invalid teacher Id")
		}

		var teacherFromDb models.Teacher
		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&teacherFromDb.ID, &teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class, &teacherFromDb.Subject)

		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				return utils.ErrorHandler(err, "Teacher not found")
			}
			return utils.ErrorHandler(err, "Error updating teachers")
		}

		teacherVal := reflect.ValueOf(&teacherFromDb).Elem()
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)
				// if field.Tag.Get("json") == k+",omitempty" {
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := teacherVal.Field(i)
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
		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Email, teacherFromDb.Class, teacherFromDb.Subject, teacherFromDb.ID)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Error updating teachers")
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

func PatchTeacherByIdDbHandler(id int, updates map[string]string) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers where id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Teacher{}, utils.ErrorHandler(err, "Teacher not found")
		}
		return models.Teacher{}, utils.ErrorHandler(err, "Error updating teacher")
	}

	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if teacherVal.Field(i).CanSet() {
					fieldVal := teacherVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Error updating teacher")
	}
	return existingTeacher, nil
}

func DeleteTeachersDbHandler(ids []int) ([]int, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error starting transaction")
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
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
			return nil, utils.ErrorHandler(err, "Teacher not found")
		} else if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
	}

	if len(deletedIds) != len(ids) {
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "Error deleting teachers")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "Error commiting deleted changes")
	}
	return deletedIds, nil
}
