package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func init() {
	log.Println("Starting the application...")
}

func main() {
	dbURL := os.Getenv("DB_HOST")

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
		return
	}

	defer dbConn.Close()

	//=================================================================================================

	log.Println("Connected to the database")

	createTab := `CREATE TABLE IF NOT EXISTS users ( id SERIAL PRIMARY KEY, name TEXT, age INT );`
	_, err = dbConn.Exec(createTab)
	if err != nil {
		log.Fatal("Error creating the table: ", err)
		return
	}
	log.Println("Created the table")

}
