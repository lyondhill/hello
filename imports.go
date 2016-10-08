package main

import "github.com/nanobox-io/golang-portal-client"
import "github.com/nanopack/portal/core"

func main() {
	s := portal.Server{
		// todo: change "Id" to "name" (for clarity)
		Id: "1",
		Host: "2",
		Port: 1,
		Forwarder: "wow",
		Weight: 100,
		UpperThreshold: 1,
		LowerThreshold: 1,
	}	
	portal.ShowServer(s)

	service := portal.Service{
		Service: core.Service{
			Id: "1",
			Host: "2",
			Port: 23,
			Type: "tcp",
			Scheduler: "what",
			Persistence: 1,
			Netmask: "apples",			
		},
		Servers: []portal.Server{s},
	}
	portal.ShowService(service)
}