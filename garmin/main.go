package main

import (
	"time"
	"fmt"
	"os"
	"github.com/abrander/garmin-connect"
)

type logger struct {

}
func (logger) Printf(format string, v ...interface{}) {
	fmt.Print("logging: ")
	fmt.Printf(format, v...)
}

func main()  {
	c := connect.NewClient(
		connect.Credentials("lyondhill@gmail.com", "GusandR2"), 
		connect.DumpWriter(os.Stdout),
		connect.DebugLogger(logger{}),
	)
	fmt.Println(c.Authenticate())

	ds, err := c.DailySummary("", time.Now())
	fmt.Println("ds", ds)
	fmt.Println("err", err)
}