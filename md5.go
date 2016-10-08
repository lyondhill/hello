package main


import (
	"crypto/md5"  
	"fmt"
)

func main() {
	hash := md5.New()
	dataall := []byte("what\n")

	hash.Write(dataall)
	fmt.Printf("%x\n", md5.Sum(dataall))

	fmt.Printf("%x %d\n", hash.Sum(nil), hash.Size())

}