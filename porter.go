package main

import (
	"os"
	"fmt"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("give me a name dummy")
	}
	b := []byte(os.Args[1])
	num := 0
	for i := 0; i < len(b); i++ {
		num = num + int(b[i])
	}
	for num < 1000 {
		num = num * 10
	}
	fmt.Println(num)
}