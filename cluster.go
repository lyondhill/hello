package main

import "fmt"
import "github.com/influxdata/plutonium/meta/control"

func main() {
	client := control.NewClient("10.0.154.249:8091")
	c, e := client.ShowCluster()
	fmt.Printf("%+v, %+v\n", c, e)
}
