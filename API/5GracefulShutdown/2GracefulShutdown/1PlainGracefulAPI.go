package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("สวัสดี")
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`death`))
	})

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`pingkunga`))
	})

	srv := http.Server{
		Addr:    ":10170",
		Handler: mux,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	fmt.Println("server starting at :10170")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	fmt.Println("shutting down...")
	if err := srv.Shutdown(context.Background()); err != nil {
		//ตรงนี้จะทำอะไรก็ทำ ทำ Request ที่รับมาแล้ว
		//บอก K8S ฝากลูกเมียข้าด้วย ให้เรียบร้อยยย
		//เฮือกกกกก
		fmt.Println("shutdown err:", err)
	}
	fmt.Println("bye bye")
}
