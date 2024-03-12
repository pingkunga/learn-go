// Sum Package
// #01 How one
//
//	Foo ba
//
// #02 How two
//
//	Ba Foo
package sum

// General Sum
// Note: Sum S = Uppper case for visible from a different package,
//
// Example:
//
// 		result := Sum(1, -2)
//
func Sum(a, b int) int {
	return a + b
}

// Sum from array input\
//
// Example:
//
//		result := SumVariadic(1, 2, 3, 4)
//
func SumVariadic(xs ...int) int {
	var total int
	//foreach ในภาษาอื่นๆ แหละ
	// _ บอก current index
	// แต่ละ Element เก็บค่าลงตัวแปร num ใน Loop แต่ละรอบ
	for _, num := range xs {
		total += num
	}
	return total
}
