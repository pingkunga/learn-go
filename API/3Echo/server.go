package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//Convert from json_http_mux_authMiddleware.go

type User struct {
	Id   int    `json:"userid"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = []User{
	{Id: 1, Name: "PingkungA", Age: 33},
}

func getUserHandler(c echo.Context) error {

	/* Old Manual Logic
	if r.Method == "GET" {
		//log.Println("GET")  	--->> <Move to Log Middleware>
		bj, errj := json.Marshal(users)

		if errj != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			w.Write([]byte(errj.Error()))
			return
		}
		w.Write(bj)
		return
	}
	*/

	return c.JSON(http.StatusOK, users)
}

type Err struct {
	Message string `json:"message"`
}

func createUserHandler(c echo.Context) error {
	/* Old Manual Logic
	if r.Method == "POST" {
		//log.Println("POST")	--->> <Move to Log Middleware>
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "error : %v", err)
			return
		}

		u := User{}
		err = json.Unmarshal(body, &u)
		if err != nil {
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		users = append(users, u)
		fmt.Printf("% #v\n", users)

		fmt.Fprintf(w, "hello %s created users", "POST")
		return
	}
	*/
	usr := User{}
	if err := c.Bind(&usr); err != nil {
		//return err

		//Format error message
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})

	}
	users = append(users, usr)
	fmt.Println("Created user:", usr)
	return c.JSON(http.StatusCreated, usr)
}

func main() {
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
	authoriedRoute := srv.Group("/")

	//Middleware-Auth
	authoriedRoute.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "admin" && password == "admin" {
			return true, nil
		}
		return false, nil
	}))

	authoriedRoute.GET("users", getUserHandler)
	authoriedRoute.POST("users", createUserHandler)

	//Note
	//Group (/api)
	//	- GET /users
	//	- POST /users

	log.Fatal(srv.Start(":10170"))
}
