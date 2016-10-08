package main

import (
  "net/http"
  "net/url"
  "fmt"
)


func main() {
  url, err := url.Parse("https://api.pagodagrid.io/organizations/")
  if err != nil {
    fmt.Println(err)
  }
  q := url.Query()
  q.Set("auth_token", "V7WfibKFZ2c8qBsBy46S")
  url.RawQuery = q.Encode()

  req, err := http.NewRequest("GET", url.String(), nil)
  if err != nil {
    fmt.Println(err)
  }
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Content-Type", "application/json")

  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(resp)

}
