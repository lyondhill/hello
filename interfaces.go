package main


import(
	"net"
	"fmt"
)

func main() {
	Connected()
}


// check to see if the bridge is connected
func Connected() bool {

	interfaces, err := net.Interfaces()
	if err != nil {
		return false
	}

	// look throught the interfaces on the system
	for _, i := range interfaces {
		fmt.Println(i)
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		// find a
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}
			fmt.Println("ip", ip)
		}
	}

	return false
}
