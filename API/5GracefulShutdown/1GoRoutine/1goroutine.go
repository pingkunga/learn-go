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
	go slow("gopher")

	slow("pingkunga")
	fmt.Println("task2 done. Sec: ", time.Since(startTime).Seconds())

	time.Sleep(10 * time.Second)
	fmt.Println("all task done. Sec: ", time.Since(startTime).Seconds())

}
