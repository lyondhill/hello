package main

import (
	"log"	
	"os"
	"time"
	"net/http"
)

func main() {
	client := http.Client{
		Timeout: 5*time.Second,
	}

	online := false
	for {
		<-time.After(time.Second)
		_, err := client.Get(os.Args[1])
		if err != nil {
			if online {
				online = false
				log.Printf("going down")
			}
			continue
			// handle error
		}
		if !online {
			online = true				
			log.Println("online")
		}
	}

}
