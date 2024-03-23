package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type User struct {
	Id   int    `json:"userid"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = []User{
	{Id: 1, Name: "PingkungA", Age: 33},
}

func userHandler(w http.ResponseWriter, r *http.Request) {
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
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("Log Middleware %s (%s) %s %s Host: %s Start Time: %s ", r.RemoteAddr, r.Proto, r.Method, r.URL, r.Host, startTime)
		next.ServeHTTP(w, r)
		log.Printf("Log Middleware %s (%s) %s %s Host: %s Duration: %d ns ", r.RemoteAddr, r.Proto, r.Method, r.URL, r.Host, time.Since(startTime).Nanoseconds())
	}
}

func main() {
	http.HandleFunc("/users", logMiddleware(userHandler))
	http.HandleFunc("/healthcheck", logMiddleware(healthCheckHandler))

	log.Println("Starting server on :10170")
	log.Fatal(http.ListenAndServe(":10170", nil))
	log.Println("Server Stopped")
}
