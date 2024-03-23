package main

import (
	"fmt"
)

type Decorator func(s string) error

func Use(next Decorator) Decorator {
	//Warpper next(s)
	return func(src string) error {
		fmt.Println("Use Decorator")
		updated := src + " decorated"
		return next(updated)
	}
}

func hello(src string) error {
	fmt.Println("Hello", src)
	return nil
}

func main() {
	WarppedFuncHome := Use(hello)
	w := WarppedFuncHome("PingkungA")
	fmt.Println("end of main", w)
}
