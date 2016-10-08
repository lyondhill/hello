package main


import (
	"bitbucket.org/nanobox/nanoauth"
	"net/http"
)

func main() {
	nanoauth.ListenAndServeTLS("0.0.0.0:2000", "token", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("hello friend"))
	}))
}