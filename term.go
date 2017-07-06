package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"bytes"
)

func main() {
	buffer := &bytes.Buffer{}
	term := terminal.NewTerminal(buffer, "")
	// fmt.Println(term.ReadLine())
	data := []byte("\x1b[2K\x1b[1G\x1b[1G       0/960\x1b[2K\x1b[4G\x1b[4G 0/17333\x1b[4G 128/17333\n")
	term.Write(data)
	// fmt.Println(term.ReadLine())
	line, _ := term.ReadLine()
	fmt.Printf("%q", line)

}
