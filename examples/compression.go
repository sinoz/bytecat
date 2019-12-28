package main

import (
	"compress/gzip"
	"fmt"
	"log"

	"github.com/sinoz/bytecat"
)

func main() {
	bldr := bytecat.NewDefaultBuilder()
	bldr.WriteInt32(8) // writing 8 bytes

	writer := gzip.NewWriter(bldr)
	if _, err := writer.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8}); err != nil {
		log.Fatal(err)
	}

	if err := writer.Flush(); err != nil {
		log.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		log.Fatal(err)
	}

	bs := bytecat.StringWrap(bldr.Build().ToByteArray())
	itr := bs.Iterator()

	uncompressedLength, _ := itr.ReadUInt32()

	reader, err := gzip.NewReader(itr)
	if err != nil {
		log.Fatal(err)
	}

	resultBytes := make([]byte, int(uncompressedLength))
	if _, err := reader.Read(resultBytes); err != nil {
		log.Fatal(err)
	}

	fmt.Println(resultBytes)
}
