package main

import (
	"hash/crc32"  
	"fmt"
)

func main() {
	table := crc32.MakeTable(crc32.IEEE)
	hash := crc32.New(table)
	dataall := []byte("helloguy\n")

	hash.Write(dataall)

	fmt.Printf("%d %d\n", hash.Sum32(), hash.Size())
}