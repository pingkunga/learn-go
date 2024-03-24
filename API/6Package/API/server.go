package main

import (
	"context"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	authoriedRoute.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "admin" && password == "admin" {
			return true, nil
		}
		return false, nil
	}))

	authoriedRoute.GET("/users", user.GetUserListHandler)
	authoriedRoute.GET("/users/:id", user.GetUserHandler)
	authoriedRoute.POST("/users", user.CreateUserHandler)
	authoriedRoute.PUT("/users/:id", user.UpdateUserHandler)
	authoriedRoute.DELETE("/users/:id", user.DeleteUserHandler)

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
