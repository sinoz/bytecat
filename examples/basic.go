package main

import (
	"fmt"

	"github.com/sinoz/bytecat"
)

func main() {
	s1 := bytecat.StringOf(1, 2, 3, 4)
	s2 := bytecat.StringOf(5, 6, 7, 8)

	fmt.Println(s1.Concat(s2)) // &{[1 2 3 4 5 6 7 8 ]}
}
