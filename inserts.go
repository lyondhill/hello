package main

import (
    "database/sql"
    _ "github.com/lib/pq"
    "fmt"
    "time"
)

func main() {

		for i := 0; i < 50; i++ {
			go func() {
		    db, err := sql.Open("postgres", "dbname=test user=postgres sslmode=disable host=75.126.15.13 port=5432")
		    checkErr(err)

				for j := 0; j < 100; j++ {

			    //Insert
			    stmt, err := db.Prepare("INSERT INTO guy(name,age) VALUES($1,$2)")
			    checkErr(err)

			    res, err := stmt.Exec("greg", j)
			    checkErr(err)
			    fmt.Println(res)
					
				}
			}()
		}
		time.Sleep(1000 * time.Second)
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}