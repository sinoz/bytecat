package main

import (
	"fmt"

	"github.com/sinoz/bytecat"
)

func main() {
	bs := bytecat.
		NewDefaultBuilder().
		WriteByte(1).
		WriteInt32(1 << 18).
		WriteInt64(1 << 35).
		Build()

	itr := bs.Iterator()
	v1, _ := itr.ReadByte()
	v2, _ := itr.ReadUInt32()
	v3, _ := itr.ReadUInt64()

	fmt.Println(v1) // 1
	fmt.Println(v2) // 262144
	fmt.Println(v3) // 34359738368
}
