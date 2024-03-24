package main

import (
	"fmt"
	"time"
)

func slow(s string) {
	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println(s, ":", i)
	}
}

func main() {
	startTime := time.Now()

	done := make(chan bool)

	//1 Done
	go func() {
		slow("gopher")
		done <- true
	}()

	slow("pingkunga")

	<-done //wait for goroutine to finish

	fmt.Println("task2 done. Sec: ", time.Since(startTime).Seconds())
	fmt.Println("all task done. Sec: ", time.Since(startTime).Seconds())

}
