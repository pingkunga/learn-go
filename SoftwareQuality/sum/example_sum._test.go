package sum

import "fmt"

func ExampleSum() {
	result := Sum(1, -2)
	fmt.Println(result)
	// Output: -1
}

func ExampleSumVariadic() {
	result := SumVariadic(1, 2, 3, 4)
	fmt.Println(result)
	// Output: 10
}
