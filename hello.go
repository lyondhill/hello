
package main

import (
	"fmt"
	"net"
	// "os"
)

func main() {
	fmt.Println(localIP())
}

func localIP() net.IP {
	tt, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, t := range tt {
		aa, err := t.Addrs()
		if err != nil {
			return nil
		}
		for _, a := range aa {
			ipnet, ok := a.(*net.IPNet)
			if !ok {
				continue
			}

			v4 := ipnet.IP.To4()
			if v4 == nil || v4[0] == 127 { // loopback address
				continue
			}
			return v4
		}
	}
	return nil
}
