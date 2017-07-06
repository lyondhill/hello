package main

import (
	"fmt"
	"os/exec"
	"os"
)

func main() {
	programName := "nanobox"
	thing, err := os.Stat(programName)

	fmt.Println("file exists:", programName, thing, err)

	// lookup the full path to nanobox
	path, err := exec.LookPath(programName)
	fmt.Println("look path results:", path, err)
	
}
