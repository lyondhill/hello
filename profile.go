package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	t := time.Now().Format("2006-01-02T15_04")
	cmd := exec.Command("curl", "-o", fmt.Sprintf("/home/ubuntu/profile-%s.tar.gz", t), "http://localhost:8086/debug/pprof/all")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(out, err)
		os.Exit(1)
	}

}
