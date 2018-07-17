package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func main() {

	go http.ListenAndServe("0.0.0.0:8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for num := 0; num < 100; num++ {
			select {
			case <-r.Context().Done():
				fmt.Println("what")
				return
			default:
			}

			time.Sleep(time.Second)
			fmt.Println("server", num)
			w.Write([]byte(strconv.Itoa(num)))
		}
	}))

	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	fmt.Println(err)

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)

	req = req.WithContext(ctx)

	go func() {
		time.Sleep(4 * time.Second)
		fmt.Println("cancel")
		cancel()
		http.DefaultTransport.(*http.Transport).CancelRequest(req)
	}()

	resp, err := http.DefaultClient.Do(req)
	fmt.Println(resp, "err", err)
	time.Sleep(time.Minute)
}
