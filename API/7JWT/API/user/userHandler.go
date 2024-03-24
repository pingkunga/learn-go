package user

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUserListHandler(c echo.Context) error {

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

func GetUserHandler(c echo.Context) error {
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

func CreateUserHandler(c echo.Context) error {
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

func UpdateUserHandler(c echo.Context) error {
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

func DeleteUserHandler(c echo.Context) error {
	id := c.Param("id")
	_, err := dbConn.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
