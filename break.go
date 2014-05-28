package main

import (
	"fmt"
)

func main() {

	fmt.Println("before loop")
magic:
	for i := 0; i < 100; i++ {
		if i < 10 {
			fmt.Println("num:", i)
		} else {
			break magic
		}
	}

	fmt.Println("After loop")

}
