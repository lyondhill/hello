package main

import (
	"fmt"
	"time"
)



func main() {
	var t time.Time
	if t.Equal(time.Time{}) {
		fmt.Println("no time")
	}
	fmt.Printf("%+v\n", t)
}