package user

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var dbConn *sql.DB

func InitDB() {
	//Connect to database
	dbConn = createConnection()
	if dbConn == nil {
		log.Fatal("Error connecting to the database")
		return
	}
	//Create table
	CreateTable(dbConn)
}

func createConnection() *sql.DB {
	dbURL := os.Getenv("DB_HOST")

	//dbConn Local Scope and New from :=
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
		return nil
	}

	return dbConn
}

func CreateTable(dbConn *sql.DB) {
	createTab := `CREATE TABLE IF NOT EXISTS users ( id SERIAL PRIMARY KEY, name TEXT, age INT );`
	_, err := dbConn.Exec(createTab)
	if err != nil {
		log.Fatal("Error creating the table: ", err)
		return
	}
}
