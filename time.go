package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now().Local()
	fmt.Println(t.Format("2006-01-02T15:04:05Z07:00"))
	fmt.Println(t.Format("2006-01-02T15:04:05"))
}
