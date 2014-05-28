package main

import (
  "bufio"
  "log"
  "os/exec"
)

func main() {
  cmd := exec.Command("puma", "-C", "/Users/lyon/pagoda/git/gritty/puma.rb")
  // cmd := exec.Command("ruby", "long.rb")
  log.Print(cmd.Args)
  stdout, _ := cmd.StdoutPipe()
  stderr, _ := cmd.StderrPipe()

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

  readererr := bufio.NewReader(stderr)
  go func() {
    for {
      str, _, err := readererr.ReadLine()
      if err != nil {
        return
      }
      log.Printf("err: %q", str)
    }
  }()

  if err := cmd.Run(); err != nil {
    log.Fatal(err)
  }

  // if err := cmd.Wait(); err != nil {
  //   log.Fatal(err)
  // }
}
