package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.Handle("/", handler{say: "stuff"}).Host("say.my.com").PathPrefix("/stuff/")
	r.Handle("/", handler{say: "woot"}).Host("my.com").PathPrefix("/woot/")
	r.Handle("/", handler{say: "wow"}).Host("i.say.my.com")
	http.ListenAndServe("0.0.0.0:8080", r)
}

type handler struct {
	say string
}

func (self handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte(self.say))
}