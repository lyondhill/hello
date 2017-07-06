package main

import "fmt"

func main() {
	for {
		for {
			fmt.Println("inside")
			break
		}
		fmt.Println("dude")
	}
	fmt.Println("Hello, playground")
}