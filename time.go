package main

import (
	"fmt"
	"time"
)

func main() {
	timer := time.NewTimer(time.Second)
	
	go func() {
		for {
			select {
			case <-timer.C:
				fmt.Println("ticking")
			}
		}
	}()

	time.Sleep(2*time.Second)
	timer.Reset(time.Second)
	timer.Stop()
	timer.Reset(time.Second)
	time.Sleep(2*time.Second)
}
