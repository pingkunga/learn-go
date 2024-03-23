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
	//stmt, err := dbConn.Prepare("SELECT id, name, age FROM users")
	stmt, err := dbConn.Prepare("SELECT id, name, age FROM users where id=$1")
	if err != nil {
		log.Fatal("Error preparing the statement: ", err)
		return
	}

	//Execute the query
	userId := 1
	rows, err := stmt.Query(userId)
	if err != nil {
		log.Fatal("Error executing the query: ", err)
		return
	}

	//View Results
	for rows.Next() {
		var id int
		var name string
		var age int

		err = rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatal("Error scanning the row: ", err)
			return
		}

		log.Println("id:", id, "name:", name, "age:", age)
	}
}
