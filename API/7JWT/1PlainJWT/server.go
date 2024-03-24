package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Err struct {
	Message string `json:"message"`
}

var users = []User{
	{Id: 1, Name: "pingkunga", Age: 33},
}

func createUserHandler(c echo.Context) error {
	u := User{}
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	users = append(users, u)

	fmt.Println("id : % #v\n", u)

	return c.JSON(http.StatusCreated, u)
}

func getUsersHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	//name := claims.Name
	fmt.Printf("User Name : %#v\n", claims)

	return c.JSON(http.StatusOK, users)
}

func AuthMiddleware(username, password string, c echo.Context) (bool, error) {
	if username == "admin" || password == "admin" {
		return true, nil
	}
	return false, nil
}
func jwtMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.ErrUnauthorized
		}
		parts := strings.Split(token, " ")
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			return echo.ErrUnauthorized
		}
		token = parts[1]

		//Secret should not hardcode
		token2, err := jwt.ParseWithClaims(token, &jwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			return echo.ErrUnauthorized
		}

		// for dump
		claims, ok := token2.Claims.(*jwtCustomClaims)
		if !ok {
			return echo.ErrUnauthorized
		}
		fmt.Printf("claims : %#v\n", claims)
		c.Set("user", token2)
		return next(c)
	}
}

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func login(c echo.Context) error {

	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username != "admin" || password != "admin" {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		"admin minda",
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 30)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	g := e.Group("/api")
	g.POST("/auth", login)

	g.Use(jwtMiddleWare)
	g.POST("/users", createUserHandler)
	g.GET("/users", getUsersHandler)

	log.Fatal(e.Start(":10170"))
}
