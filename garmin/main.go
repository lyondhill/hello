package main

import (
	"fmt"
	"os"
	"time"

	connect "github.com/abrander/garmin-connect"
)

type logger struct {
}

func (logger) Printf(format string, v ...interface{}) {
	fmt.Print("logging: ")
	fmt.Printf(format, v...)
}

func main() {
	username, _ := os.LookupEnv("USERNAME")
	password, _ := os.LookupEnv("PASSWORD")
	c := connect.NewClient(
		connect.Credentials(username, password),
		connect.DumpWriter(os.Stdout),
		connect.DebugLogger(logger{}),
	)
	fmt.Println(c.Authenticate())

	ds, err := c.DailySummary("", time.Now())
	fmt.Println("ds", ds)
	fmt.Println("err", err)
}
