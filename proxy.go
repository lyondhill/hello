package main

import (
  "flag"
  "log"
  "net/http"
  "io/ioutil"
  // "strings"
  // "unicode"
)

var (
  listen = flag.String("listen", "localhost:1080", "listen on address")
  logp = flag.Bool("log", false, "enable logging")
)

func main() {
  flag.Parse()
  proxyHandler := http.HandlerFunc(report)
  log.Fatal(http.ListenAndServe(*listen, proxyHandler))

}


func report(w http.ResponseWriter, r *http.Request){

  uri := "http://drawception.com"+r.RequestURI

  log.Println(r.Method + ": " + uri)

  if r.Method == "POST" {
    body, err := ioutil.ReadAll(r.Body)
    fatal(err)
    log.Printf("Body: %v\n", string(body));
  }

  rr, err := http.NewRequest(r.Method, uri, r.Body)
  fatal(err)
  copyHeader(r.Header, &rr.Header)

  // Create a client and query the target
  var transport http.Transport
  resp, err := transport.RoundTrip(rr)
  fatal(err)

  log.Printf("Resp-Headers: %v\n", resp.Header);

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  fatal(err)

  dH := w.Header()
  copyHeader(resp.Header, &dH)
  dH.Add("Requested-Host", rr.Host)

  w.Write(body)
}

func fatal(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func copyHeader(source http.Header, dest *http.Header){
  for n, v := range source {
      for _, vv := range v {
          dest.Add(n, vv)
      }
  }
}