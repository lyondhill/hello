package main

import "fmt"
import "net"
import "io"

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:5555")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return		
	}

	io.Copy(conn, conn)
	conn.Close()
	listener.Close()
}