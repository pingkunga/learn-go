package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

//Convert from json_http_mux_authMiddleware.go

type User struct {
	Id   int    `json:"userid"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Err struct {
	Message string `json:"message"`
}

func getUserListHandler(c echo.Context) error {

	//Create Prepared Statement
	stmt, err := dbConn.Prepare("SELECT id, name, age FROM users")
	if err != nil {
		log.Fatal("Error preparing the statement: ", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query all users statment:" + err.Error()})
	}

	//Execute the query
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("Error executing the query: ", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't query all users:" + err.Error()})
	}

	//Add result to struct array
	users := []User{}

	for rows.Next() {
		u := User{}
		//& send reference to set value
		//order same as select statement
		err := rows.Scan(&u.Id, &u.Name, &u.Age)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan user:" + err.Error()})
		}
		users = append(users, u)
	}
	return c.JSON(http.StatusOK, users)

}

func getUserHandler(c echo.Context) error {
	id := c.Param("id")
	//Create Prepared Statement
	stmt, err := dbConn.Prepare("SELECT id, name, age FROM users WHERE id = $1")
	if err != nil {
		log.Fatal("Error preparing the statement: ", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query user statment:" + err.Error()})
	}

	//Execute the query
	row := stmt.QueryRow(id)
	u := User{}
	err = row.Scan(&u.Id, &u.Name, &u.Age)
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "user not found"})
	case nil:
		return c.JSON(http.StatusOK, u)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan user:" + err.Error()})
	}

}

func createUserHandler(c echo.Context) error {
	/*
		Old Code
		usr := User{}
		if err := c.Bind(&usr); err != nil {
			//return err

			//Format error message
			return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})

		}
		users = append(users, usr)
		fmt.Println("Created user:", usr)
		return c.JSON(http.StatusCreated, usr)
	*/

	u := User{}
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := dbConn.QueryRow("INSERT INTO users (name, age) values ($1, $2)  RETURNING id", u.Name, u.Age)
	err = row.Scan(&u.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, u)
}

func updateUserHandler(c echo.Context) error {
	id := c.Param("id")
	u := User{}
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	_, err = dbConn.Exec("UPDATE users SET name=$1, age=$2 WHERE id=$3", u.Name, u.Age, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, u)
}

func deleteUserHandler(c echo.Context) error {
	id := c.Param("id")
	_, err := dbConn.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

var dbConn *sql.DB

func main() {
	//Connect to database
	dbConn = createConnection()
	defer dbConn.Close()

	//Create table
	CreateTable(dbConn)

	srv := echo.New()

	//Middleware-Log
	srv.Use(middleware.Logger())
	srv.Use(middleware.Recover())

	//Unauthorized
	//Minamal API Like DotNet
	srv.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	//Authorized
	authoriedRoute := srv.Group("/api")

	//Middleware-Auth
	authoriedRoute.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "admin" && password == "admin" {
			return true, nil
		}
		return false, nil
	}))

	authoriedRoute.GET("/users", getUserListHandler)
	authoriedRoute.GET("/users/:id", getUserHandler)
	authoriedRoute.POST("/users", createUserHandler)
	authoriedRoute.PUT("/users/:id", updateUserHandler)
	authoriedRoute.DELETE("/users/:id", deleteUserHandler)

	log.Fatal(srv.Start(":10170"))
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
	log.Println("Created the table")
}
