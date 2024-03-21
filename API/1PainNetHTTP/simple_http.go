package main

import (
	"log"
	"net/http"
)

func handleAboutMe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{name: "PingkungA", age: 33}`))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))

	})

	http.HandleFunc("/aboutme", handleAboutMe)

	//Filter Specfic Http Method
	http.HandleFunc("/aboutme1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		w.Write([]byte(`{name: "PingkungA_GET", age: 33}`))
	})

	log.Println("Starting server on :10170")
	log.Fatal(http.ListenAndServe(":10170", nil))
	log.Println("Server Stopped")
}
