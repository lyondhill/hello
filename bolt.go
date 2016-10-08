package main

import (
    "fmt"
    "log"
    // "time"
    // "strconv"

    "github.com/boltdb/bolt"
)

var world = []byte("stats")

func main() {
    db, err := bolt.Open("./bolt.db", 0644, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // key := []byte("stats")

    db.Update(func(tx *bolt.Tx) error {
      bucket, err := tx.CreateBucketIfNotExists(world)
      if err != nil {
        return err
      }
      fmt.Printf("%+v\n", bucket.Stats())
      max := 10
      if bucket.Stats().KeyN > max {
        delete_count := bucket.Stats().KeyN - max
        fmt.Println(delete_count)
        c := bucket.Cursor()
        for i := 0; i < delete_count; i++ {
          k, v := c.First()
          fmt.Printf("%s - %s", k, v)
          fmt.Println(c.Delete())
        }
      }
      return nil

    })

    // store some data
    // err = db.Update(func(tx *bolt.Tx) error {
    //     bucket, err := tx.CreateBucketIfNotExists(world)
    //     if err != nil {
    //         return err
    //     }
    //     for i := 0; i < 100; i++ {
    //       err = bucket.Put([]byte(time.Now().String()), []byte("here is some number: "+strconv.Itoa(i)))
    //       if err != nil {
    //           return err
    //       }
    //     }
    //     return nil
    // })

    // if err != nil {
    //     log.Fatal(err)
    // }


    // retrieve the data
    err = db.View(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(world)
        if bucket == nil {
            return fmt.Errorf("Bucket %q not found!", world)
        }

        c := bucket.Cursor()
        for k, v := c.First(); k != nil; k, v = c.Next() {
            fmt.Printf("%s is %s.\n", k, v)
        }

        // val := bucket.Get(key)
        // fmt.Println(string(val))

        return nil
    })

    if err != nil {
        log.Fatal(err)
    }
}
