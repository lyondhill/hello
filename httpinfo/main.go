package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		b, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("body: %s\n", b)
		w.Write(b)
	}))
}
