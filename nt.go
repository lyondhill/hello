package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	strNum, _ := strconv.ParseInt(os.Args[1], 10, 64)
	fmt.Println(time.Unix(0, strNum))
}
