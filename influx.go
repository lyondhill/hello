package main

import (
	"strings"
	"github.com/influxdata/influxdb/client/v2"
	"fmt"
	"time"
	"github.com/Pallinder/go-randomdata"
)

func main() {
	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
	    Addr: "http://localhost:8086",
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	r, err := c.Query(client.NewQuery("CREATE DATABASE statistics", "statistics", "s"))
	if err != nil {
		panic(err)
	}
	fmt.Println(r)

	r, err = c.Query(client.NewQuery(`CREATE RETENTION POLICY "5.hour" ON statistics DURATION 5h REPLICATION 1 DEFAULT`, "statistics", "s"))
	if err != nil {
		panic(err)
	}
	fmt.Println(r)

	r, err = c.Query(client.NewQuery(`CREATE RETENTION POLICY "1.hour" ON statistics DURATION 1h REPLICATION 1 DEFAULT`, "statistics", "s"))
	if err != nil {
		panic(err)
	}
	fmt.Println(r)

	go randomInsert(c)
	keepContinuousQueriesUpToDate(c)
	
	
}
func keepContinuousQueriesUpToDate(c client.Client) {
	for {
		cols, err := c.Query(client.NewQuery("select * from \"1.hour\".metrics limit 1", "statistics", "s"))
		if err != nil	{
			panic(err)
		}

		// populate current columns
		columns := []string{}
		for _, res := range cols.Results {
			for _, series := range res.Series {
				if series.Name == "metrics" {
					for _, col := range series.Columns {
						if col != "time" && col != "cpu" {
							str := fmt.Sprintf("mean(%s) as %s", col, col)
							columns = append(columns, str)
						}
					}
				}
			}
		}

		colString := strings.Join(columns, ", ")

		r, err := c.Query(client.NewQuery(`CREATE CONTINUOUS QUERY "aggrigate" ON statistics BEGIN select `+colString+` into "5.hour"."metrics" from "1.hour"."metrics" group by time(5m), cpu END`, "statistics", "s"))
		if err != nil {
			fmt.Printf("ERROR: %+v, %+v\n", r, err)
		}
		fmt.Println("adding continuous query for", colString)


		<-time.After(time.Minute)
	}
}

// func keepContinuousQueriesUpToDate(c client.Client) {
// 	for  {
// 		cols, err := c.Query(client.NewQuery("select * from \"1.hour\".metrics limit 1", "statistics", "s"))
// 		if err != nil	{
// 			panic(err)
// 		}
// 		cont, err := c.Query(client.NewQuery("SHOW CONTINUOUS QUERIES", "statistics", "s"))
// 		if err != nil	{
// 			panic(err)
// 		}


// 		// populate current columns
// 		columns := []string{}
// 		for _, res := range cols.Results {
// 			for _, series := range res.Series {
// 				if series.Name == "metrics" {
// 					columns = append(columns, series.Columns...)
// 				}
// 			}
// 		}

// 		// populate current continuous queries
// 		cqs := []string{}
// 		for _, res := range cont.Results {
// 			for _, series := range res.Series {
// 				if series.Name == "statistics" {
// 					for _, val := range series.Values {
// 						st, ok := val[0].(string)
// 						if ok {
// 							cqs = append(cqs, st)
// 						}
// 					}
// 				}
// 			}
// 		}

// 		fmt.Println("columns:", columns)
// 		fmt.Println("cqs:", cqs)

// 		for _, col := range columns {
// 			if col != "cpu" && col != "time" {
// 				add := true
// 				for _, cq := range cqs {
// 					if cq == col {
// 						add = false
// 					}
// 				}
// 				if add {
// 					r, err := c.Query(client.NewQuery(`CREATE CONTINUOUS QUERY "`+col+`" ON statistics BEGIN select mean(`+col+`) as `+col+` into "5.hour"."metrics" from "1.hour"."metrics" group by time(5m), cpu END`, "statistics", "s"))
// 					if err != nil {
// 						fmt.Printf("ERROR: %+v, %+v\n", r, err)
// 					}
// 					fmt.Println("adding continuous query for", col)
// 				}
// 			}
// 		}
		


// 		for _, cq := range cqs {
// 			remove := true
// 			for _, col := range columns {
// 				if col == cq {
// 					remove = false
// 				}
// 			}
// 			if remove {
// 				r, err := c.Query(client.NewQuery(`DROP CONTINUOUS QUERY "`+cq+`" ON statistics`, "statistics", "s"))
// 				if err != nil {
// 					fmt.Printf("ERROR: %+v, %+v\n", r, err)
// 				}
// 				fmt.Println("removing continuous query for", cq)
// 			}
// 		}

// 		<-time.After(time.Minute)	
// 	}


// }

func randomInsert(c client.Client) {
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
	    Database:  "statistics",
	    Precision: "s",
	    RetentionPolicy: "1.hour",
	})

	cpus := []string{randomdata.FirstName(randomdata.Female), randomdata.FirstName(randomdata.Female), randomdata.FirstName(randomdata.Female)}
	fmt.Println("new computers:", cpus)

	statName := []string{randomdata.SillyName(),randomdata.SillyName(),randomdata.SillyName(),randomdata.SillyName()}
	fmt.Println("new stats:", statName)

	for {
		// Create a point and add to batch
		for _, cpu := range cpus {
			tags := map[string]string{"cpu": cpu}
			fields := map[string]interface{}{}
			for _, stat := range statName {
				fields[stat] = randomdata.Decimal(0, 100)
			}
			pt, err := client.NewPoint("metrics", tags, fields, time.Now())
			if err != nil {
			    fmt.Println("Error: ", err.Error())
			}
			bp.AddPoint(pt)
		}
		// Write the batch
		c.Write(bp)	

		<-time.After(time.Minute)
	}
	
}
