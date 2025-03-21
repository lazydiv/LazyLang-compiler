package main

import "fmt"

func main() {
	nums := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	x := 0
	for ; (x < 9); x = (x + 1) {
	if (nums[x] == 3) {
		fmt.Println(x)
	}
}
}
