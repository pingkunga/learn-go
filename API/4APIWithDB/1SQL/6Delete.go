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

	//Delete
	//Create Prepared Statement
	stmt, err := dbConn.Prepare("DELETE FROM users WHERE id=$1;")
	if err != nil {
		log.Fatal("Error preparing the statement: ", err)
		return
	}

	//Execute the query
	userId := 2
	_, err = stmt.Exec(userId)
	if err != nil {
		log.Fatal("Error executing the query: ", err)
		return
	}

	log.Println("Delete successful")
}
