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

func GetExecByIdDbHandler(id int) (models.Exec, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Error connecting to DB")
	}
	defer db.Close()

	var exec models.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username class FROM execs WHERE id = ?", id).Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username)

	if err == sql.ErrNoRows {
		return models.Exec{}, utils.ErrorHandler(err, "exec Not found")
	} else if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Error retrieving exec")
	}
	return exec, nil
}

func GetExecsDbHandler(execs []models.Exec, r *http.Request) ([]models.Exec, error) {

	query := "SELECT id, first_name, last_name, email, username FROM execs WHERE 1 = 1"
	var args []interface{}

	query, args = utils.AddFilters(r, query, args)

	// also handling - execs/?sortby=name:asc&sortby=class:desc
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
		var exec models.Exec
		err = rows.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error scaning row from db")
		}
		execs = append(execs, exec)
	}
	return execs, nil
}

func AddExecsDbHandler(addedExecs []models.Exec, newExecs []models.Exec) ([]models.Exec, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	// stmt, err := db.Prepare("INSERT INTO Execs (first_name, last_name, email, class) VALUES(?,?,?,?,?)")
	stmt, err := db.Prepare(utils.GenerateInsertQuery("execs", models.Exec{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error updating Execs")
	}
	defer stmt.Close()

	for _, newExec := range newExecs {
		values := utils.GetStructValues(newExec)
		res, err := stmt.Exec(values...)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				return nil, utils.ErrorHandler(err, "")
			}
			return nil, utils.ErrorHandler(err, "Error updating Exec")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error updating Exec")
		}
		newExec.ID = int(lastID)
		addedExecs = append(addedExecs, newExec)
	}
	return addedExecs, nil
}

func DeleteExecByIdDbHandler(id int) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM execs WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "Error deleting Exec")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "Error deleting Exec")
	}

	if rowsAffected == 0 {
		return utils.ErrorHandler(err, "Exec not found")
	}
	return nil
}

func PatchExecsDbHandler(updates []map[string]interface{}) error {
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
			return utils.ErrorHandler(err, "Invalid Exec Id")
		}

		var execFromDB models.Exec
		err = db.QueryRow("SELECT id, first_name, last_name, username, role FROM execs WHERE id = ?", id).Scan(&execFromDB.ID, &execFromDB.FirstName, &execFromDB.LastName, &execFromDB.Email, &execFromDB.Username, &execFromDB.Role)

		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				return utils.ErrorHandler(err, "Exec not found")
			}
			return utils.ErrorHandler(err, "Error updating Execss")
		}

		execVal := reflect.ValueOf(&execFromDB).Elem()
		execType := execVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < execVal.NumField(); i++ {
				field := execType.Field(i)
				// if field.Tag.Get("json") == k+",omitempty" {
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := execVal.Field(i)
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
		_, err = tx.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?", execFromDB.FirstName, execFromDB.LastName, execFromDB.Email, execFromDB.Username, execFromDB.Role, execFromDB.ID)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Error updating Execss")
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

func PatchExecByIdDbHandler(id int, updates map[string]string) (models.Exec, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Error Connecting to DB")
	}
	defer db.Close()

	var existingExec models.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username, role FROM execs where id = ?", id).Scan(&existingExec.ID, &existingExec.FirstName, &existingExec.LastName, &existingExec.Email, &existingExec.Username, &existingExec.Role, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Exec{}, utils.ErrorHandler(err, "Exec not found")
		}
		return models.Exec{}, utils.ErrorHandler(err, "Error updating Exec")
	}

	execVal := reflect.ValueOf(&existingExec).Elem()
	execType := execVal.Type()

	for k, v := range updates {
		for i := 0; i < execVal.NumField(); i++ {
			field := execType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if execVal.Field(i).CanSet() {
					fieldVal := execVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(execVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ?, role = ? WHERE id = ?", existingExec.FirstName, existingExec.LastName, existingExec.Email, existingExec.Username, existingExec.Role, existingExec.ID)
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Error updating exec")
	}
	return existingExec, nil
}

// func UpdateExecByIdDbHandle(id int, updatedExec models.Exec) (models.Exec, error) {
// 	db, err := ConnectDB()
// 	if err != nil {
// 		return models.Exec{}, utils.ErrorHandler(err, "Error Connecting to DB")
// 	}
// 	defer db.Close()

// 	var existingExec models.Exec
// 	err = db.QueryRow("SELECT id, first_name, last_name, email, username, role FROM execs WHERE id = ?", id).Scan(&existingExec.ID, &existingExec.FirstName, &existingExec.LastName, &existingExec.Email, &existingExec.Username, &existingExec.Role, id)
// 	if err == sql.ErrNoRows {
// 		return models.Exec{}, utils.ErrorHandler(err, "Exec not found")
// 	} else if err != nil {
// 		return models.Exec{}, utils.ErrorHandler(err, "Error updating Exec")
// 	}

// 	updatedExec.ID = existingExec.ID
// 	_, err = db.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ?, role = ? WHERE id = ?", updatedExec.FirstName, updatedExec.LastName, updatedExec.Email, updatedExec.Username, updatedExec.Role, updatedExec.ID)
// 	if err != nil {
// 		return models.Exec{}, utils.ErrorHandler(err, "Error updating Exec")
// 	}
// 	return updatedExec, nil
// }
