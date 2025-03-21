package main

import "fmt"

func main() {
	x := 15
	y := 10
	nums := []interface{}{1, 2, 3}
	fmt.Println(nums)
	fmt.Println(x)
	if (x >= y) {
		fmt.Println(x)
	} else {
		fmt.Println(y)
	}
	a := 5
	b := 5
	if (nums[0] == 1) {
		fmt.Println(1)
	} else {
		fmt.Println(2)
	}
	if (a < b) {
		fmt.Println(100)
	} else {
		if (a > b) {
		fmt.Println(200)
	} else {
		fmt.Println(300)
	}
	}
}
