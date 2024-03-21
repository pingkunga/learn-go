package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type User struct {
	Id   int    `json:"userid"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = []User{
	{Id: 1, Name: "PingkungA", Age: 33},
}

func main() {
	TryJsonMarshal()
	TryJsonUnmarshall()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			log.Println("GET")
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
			log.Println("POST")
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
	})

	log.Println("Starting server on :10170")
	log.Fatal(http.ListenAndServe(":10170", nil))
	log.Println("Server Stopped")
}

func TryJsonMarshal() {
	u := User{Id: 1, Name: "PingkungA", Age: 33}

	b, err := json.Marshal(u)

	fmt.Println(string(b))
	fmt.Println()
	fmt.Println("byte: %s \n", b)
	fmt.Println("Error: %s \n", err)

	uls := []User{
		{Id: 1, Name: "PingkungA", Age: 33},
		{Id: 2, Name: "PingkungB", Age: 34},
	}

	bj, errj := json.Marshal(uls)

	fmt.Println(string(bj))
	fmt.Println()
	fmt.Println("byte: %s \n", bj)
	fmt.Println("Error: %s \n", errj)
}

func TryJsonUnmarshall() {
	b := []byte(`{"userid": 1, "name": "PingkungA", "age": 33}`)

	var u User
	err := json.Unmarshal(b, &u)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(u)
}
