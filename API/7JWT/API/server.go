package main

import (
	"context"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pingkunga/learn-go/user"
)

func main() {
	user.InitDB()

	srv := echo.New()
	//srv.Logger.SetLevel(log.INFO)

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
	// authoriedRoute.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	if username == "admin" && password == "admin" {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))

	authoriedRoute.POST("/auth", authHandler)
	authoriedRoute.POST("/auth/refresh", refreshTokenHandler)

	authoriedRouteUsr := authoriedRoute.Group("/users")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}
	authoriedRouteUsr.Use(echojwt.WithConfig(config))
	authoriedRouteUsr.Use(adminMiddleware("admin"))

	authoriedRouteUsr.GET("", user.GetUserListHandler)
	authoriedRouteUsr.GET("/:id", user.GetUserHandler)
	authoriedRouteUsr.POST("", user.CreateUserHandler)
	authoriedRouteUsr.PUT("/:id", user.UpdateUserHandler)
	authoriedRouteUsr.DELETE("/:id", user.DeleteUserHandler)

	//===========================================================================================
	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// defer stop()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	//Start server
	go func() {
		if err := srv.Start(":10170"); err != nil && err != http.ErrServerClosed { // Start server
			srv.Logger.Fatal("shutting down the server")
		}
	}()
	<-shutdown

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		srv.Logger.Fatal(err)
	}
}

//===========================================================================================
//JWT

type jwtCustomClaims struct {
	Name string `json:"name"`
	//Admin bool   `json:"admin"`
	Role string `json:"role"`
	Type string `json:"type"` // Added to distinguish between access and refresh tokens
	jwt.RegisteredClaims
}

func authHandler(c echo.Context) error {

	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username != "admin" || password != "admin" {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		"admin minda",
		"admin",
		"access",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}

	// Create token with claims
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	accessTokenString, err := accessToken.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	// Create refresh token
	refreshTokenClaims := &jwtCustomClaims{
		"admin minda",
		"admin",
		"refresh",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Longer-lived
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
	})
}

// refreshToken function
func refreshTokenHandler(c echo.Context) error {
	refreshTokenString := c.FormValue("refresh_token")

	//"secret" should be in the environment variable not hardcoded
	jwtSecretKey := []byte("secret")

	// Parse the token
	token, err := jwt.ParseWithClaims(refreshTokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		return echo.ErrUnauthorized
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid || claims.Type != "refresh" {
		return echo.ErrUnauthorized
	}

	// Create new access token
	newAccessTokenClaims := &jwtCustomClaims{
		claims.Name,
		claims.Role,
		"access",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessTokenClaims)
	newAccessTokenString, err := newAccessToken.SignedString(jwtSecretKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  newAccessTokenString,
		"refresh_token": refreshTokenString,
	})
}

// func roleCheckMiddleware(requiredRole string) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			user := c.Get("user").(*jwt.Token)
// 			claims := user.Claims.(*jwtCustomClaims)
// 			//!= requiredRole
// 			fmt.Printf("User Name : %#v\n", claims)
// 			fmt.Println(claims.Role)
// 			if claims.Role != requiredRole {
// 				return echo.ErrUnauthorized
// 			}

// 			return next(c)
// 		}
// 	}
// }

func adminMiddleware(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*jwtCustomClaims)
			fmt.Printf("\n\n\nUser Name : %#v\n\n\n", claims)
			if claims.Role != requiredRole {
				return echo.ErrUnauthorized
			}

			return next(c)
		}
	}
}

//===========================================================================================
