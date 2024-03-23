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
	row := dbConn.QueryRow("INSERT INTO users (name, age) values ($1, $2)  RETURNING id", "PingkungA", 33)

	//Get the id of the inserted row
	var id int
	err = row.Scan(&id)

	if err != nil {
		log.Fatal("Error inserting the row: ", err)
		return
	}
	log.Println("Inserted row with id:", id)
}
