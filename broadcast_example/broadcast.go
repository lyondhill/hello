package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	laddr, err := net.ResolveUDPAddr("udp", ":0")
	mcaddr, err := net.ResolveUDPAddr("udp", "224.0.1.60:1888")
	check(err)
	conn, err := net.ListenMulticastUDP("udp", nil, mcaddr)
	lconn, err := net.ListenUDP("udp", laddr)
	check(err)
	reader := bufio.NewReader(os.Stdin)
	go listen(conn)
	for {
		fmt.Print("Input: ")
		txt, _, err := reader.ReadLine()
		b := make([]byte, 256)
		copy(b, txt)
		check(err)
		_, err = lconn.WriteToUDP(b, mcaddr)
		check(err)
	}
}

func listen(conn *net.UDPConn) {
	for {
		b := make([]byte, 256)
		_, _, err := conn.ReadFromUDP(b)
		check(err)
		fmt.Println("read", string(b))
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
