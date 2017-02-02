package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	path, err := exec.LookPath("bash")
	if err != nil {
		log.Fatal("nobash")
	}
	fmt.Printf("bash is at '%s'\n", path)
}
