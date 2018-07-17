package main

import (
	"fmt"

	"github.com/influxdata/influxdb/client/v2"
)

func main() {

	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: "http://localhost:8086"})
	if err != nil {
		panic(err)
	}

	c.Query(client.NewQuery("CREATE DATABASE test", "", ""))
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{Database: "test"})
	fmt.Println("err", err)
	p1, err := client.NewPoint("dude", map[string]string{}, map[string]interface{}{"hi": 1.0})
	p2, err := client.NewPoint("man", map[string]string{}, map[string]interface{}{"hi": 1.0})
	p3, err := client.NewPoint("cat", map[string]string{}, map[string]interface{}{"hi": 1.0})
	p4, err := client.NewPoint("bark", map[string]string{}, map[string]interface{}{"hi": 1.0})
	fmt.Println("err", err)

	bp.AddPoint(p1)
	bp.AddPoint(p2)
	bp.AddPoint(p3)
	bp.AddPoint(p4)
	c.Write(bp)

	resp, err := c.Query(client.NewQuery("show measurements", "test", ""))
	fmt.Printf("%#v, %s\n", resp, err)

	for _, result := range resp.Results {
		for _, series := range result.Series {
			for _, value := range series.Values {
				val := value[0].(string)
				fmt.Println(val)
			}
		}
	}
}
