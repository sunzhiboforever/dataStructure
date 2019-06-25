package main

import (
	"fmt"
)

func main() {
	a := make([]int, 1)
	func(s *[]int) {
		*s = append(*s, 1)
	}(&a)
	fmt.Println(a)
	return

}


