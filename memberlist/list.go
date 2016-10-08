package main

import "fmt"
import "github.com/hashicorp/memberlist"
import "time"
import "flag"


type EventThing struct {
	Count int
}

    // NotifyJoin is invoked when a node is detected to have joined.
    // The Node argument must not be modified.
func (e EventThing) NotifyJoin(n *memberlist.Node) {
	fmt.Printf("NotifyJoin:%+v", n)
}

    // NotifyLeave is invoked when a node is detected to have left.
    // The Node argument must not be modified.
func (e EventThing) NotifyLeave(n *memberlist.Node) {
	fmt.Printf("NotifyLeave:%+v", n)
}
func (e EventThing) NotifyConflict(existing, other *memberlist.Node) {
	fmt.Println("there is a conflict", existing, other)
}

    // NotifyUpdate is invoked when a node is detected to have
    // updated, usually involving the meta data. The Node argument
    // must not be modified.
func (e EventThing) NotifyUpdate(n *memberlist.Node) {
	fmt.Printf("NotifyUpdate:%+v", n)
}

var port int
var pool string

func init() {
	flag.IntVar(&port, "port", 1234, "Port to listen on")	
	flag.StringVar(&pool, "pool", "127.0.0.1:1234", "pool")
	flag.Parse()
}

func main() {
	config := memberlist.DefaultLANConfig()
	config.Name = "hay dude3"
	config.Events = EventThing{0}
	config.Conflict = EventThing{0}
	config.BindPort = port
	config.AdvertisePort = port
	fmt.Printf("%+v\n\n", config)
	list, err := memberlist.Create(config)
	if err != nil {
	    panic("Failed to create memberlist: " + err.Error())
	}

	fmt.Printf("%+v\n\n", list)

	n, err := list.Join([]string{pool})
	if err != nil {
	    panic("Failed to join cluster: " + err.Error())
	}
	fmt.Print(n)
	fmt.Printf("%+v\n\n", config)
	// Ask for members of the cluster
	for {
		for _, member := range list.Members() {
		    fmt.Printf("Member: %s %d\n", member.Name, member.Port)
		}
		time.Sleep(time.Second)
	}

}
