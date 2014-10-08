package main

import (
  "github.com/pagodabox-api-clients/go"
  "fmt"
)


func main() {

  client := pagodabox.Client{AuthToken: "t5nuYXqs5zSv_zaokDKz"}

  // pass 'nil' when no additional options are desired
  apps, err := client.GetApps()
  if err != nil {
    panic(err)
  }
  for _, app  := range apps {
    fmt.Printf("%+v\n", app)
    // fmt.Println(app.Name)
  }
  
}
