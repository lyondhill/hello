package main

import (
  "bufio"
  "log"
  "os/exec"
)

func main() {
  cmd := exec.Command("ruby", "long.rb")
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    log.Fatal(err)
  }
  reader := bufio.NewReader(stdout)
  go func() {
    for {
      str, _, err := reader.ReadLine()
      if err != nil {
        return
      }
      log.Printf("streaming: %q", str)
    }
  }()

  if err := cmd.Start(); err != nil {
    log.Print(err)
  }

  if err := cmd.Wait(); err != nil {
    log.Print(err)
  }
}
