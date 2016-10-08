package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// Listen on TCP port 2000 on all interfaces.
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			for i := 0; i < 1000000; i++ {
				c.Write([]byte("morestuff"))
				fmt.Println("write", i)
			}

			// Shut down the connection.
			c.Close()
		}(conn)
	}
}
