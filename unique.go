package main

import (
  "fmt"
  "net"
  "crypto/md5"
)

func main() {
	fmt.Println(UniqueID())
}

func UniqueID() string {
  interfaces, _ := net.Interfaces()

  for _, i := range interfaces {
  	if i.Flags&net.FlagLoopback == net.FlagLoopback {
  		return fmt.Sprintf("%x", md5.Sum(i.HardwareAddr))
  	}
  }
  return "bad"
}
