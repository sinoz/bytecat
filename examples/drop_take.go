package main

import (
	"fmt"

	"github.com/sinoz/bytecat"
)

func main() {
	bytes := bytecat.StringOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	droppedFour := bytes.Drop(4)
	takenSix := bytes.Take(6)

	fmt.Println(droppedFour) // ${[5 6 7 8 9 10]}
	fmt.Println(takenSix)    // ${[1 2 3 4 5 6]}
}
