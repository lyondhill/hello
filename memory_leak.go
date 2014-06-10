package main

import (
  "fmt"
  "time"
  "runtime"
)

func main() {
  time.Sleep(20 * time.Second)
  for i := 0; i < 10; i++ {
    fmt.Println("getting memory")
    tmp := make([]uint32, 100000000)
    for kk, _ := range tmp {
      tmp[kk] = 0
    }
    time.Sleep(5 * time.Second)
    fmt.Println("returning memory")
    time.Sleep(5 * time.Second)
  }
  runtime.GC()
  time.Sleep(20 * time.Minute)
}
