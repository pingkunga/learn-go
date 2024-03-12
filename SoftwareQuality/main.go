package main

import (
	"fmt"

	"github.com/pingkunga/learn-go/sum" // Import the package that contains the sum function
)

func main() {
	result := sum.Sum(5, 3) // Change sum.sum to sum.Sum
	fmt.Println("Sum:", result)
}
