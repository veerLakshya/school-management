package sqlconnect

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	log.Println("Trying to connet to MariaDB...")

	connectionString := os.Getenv("CONNECTION_STRING")

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to MariaDB.")
	return db, nil
}
