package main

import (
  "github.com/tarm/goserial"
  "log"
  "os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("you need to pass in the COM and the thing you want to write to the socket")
	}
	c := &serial.Config{Name: os.Args[1], Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
    log.Fatal("connection: ", err.Error())
	}
	log.Print("connection esablished!")

	n, err := s.Write([]byte(os.Args[2]+"\n"))
	if err != nil {
	  log.Fatal("write: ",err.Error())
	}
	log.Print("message written ("+os.Args[2]+")")

	for {
		buf := make([]byte, 128)
		n, err = s.Read(buf)
		if err != nil {
		  log.Fatal("read: ", err.Error())
		}
		log.Println(n)
		log.Printf("response: %s", string(buf[:n]))
	}
}