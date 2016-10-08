package main

import (
	"time"
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		<-time.After(1 * time.Millisecond)
		// Wait for a connection.
		b := make([]byte, 10)
		_, err := conn.Read(b)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("b:%s\n", b )
	}
}
