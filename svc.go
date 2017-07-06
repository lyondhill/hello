package main


import (
	"fmt"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func main() {
	fmt.Println("here we go")
	m, err := mgr.Connect()
	fmt.Printf("%#v, %s\n", m, err)
	if err != nil {
		fmt.Printf("%#v\n",false)
		return
	}
	defer m.Disconnect()

	// check to see if we need to create at all
	s, err := m.OpenService("nanobox-server")
	fmt.Printf("%#v, %s\n", s, err)
	if err != nil {
		// jobs done
		fmt.Printf("%#v\n",false)
		return
	}
	defer s.Close()

	status, err := s.Query()
	fmt.Printf("%#v, %s\n", status, err)
	if err != nil {
		fmt.Printf("%#v\n",false)
		return
	}

	fmt.Println(status.State == svc.Running)
}
