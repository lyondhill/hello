package main

import (
    "fmt"
    "io/ioutil"
    "os"
)

func main() {
  files, _ := ioutil.ReadDir("./")
  for _, f := range files {
    if f.IsDir() {
      fmt.Println("dir:", f.Name())
    } else {
      fmt.Println("name:", f.Name(), "size:", f.Size(), "bytes")
    }
  }
  f, _ := os.Stat("./exec.go")
  fmt.Println("name:", f.Name())
  fmt.Println("size:", f.Size(), "bytes")
  fmt.Println("mode:", f.Mode())
  fmt.Println("time:", f.ModTime())

}