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

	//Create Prepared Statement
	stmt, err := dbConn.Prepare("UPDATE users SET name=$2 WHERE id=$1;")
	if err != nil {
		log.Fatal("Error preparing the statement: ", err)
		return
	}

	//Execute the query
	userId := 2
	userName := "Faii"
	_, err = stmt.Exec(userId, userName)
	if err != nil {
		log.Fatal("Error executing the query: ", err)
		return
	}

	log.Println("Update successful")
}
