package main

import (
	// "bytes"
	"io"
	"fmt"

	"github.com/docker/docker/pkg/term"
)


func main() {
	// establish file descriptors for std streams
	stdin, stdout, _ := term.StdStreams()
	stdInFD, _ := term.GetFdInfo(stdin)
	stdOutFD, _ := term.GetFdInfo(stdout)

	oldInState, err := term.SetRawTerminal(stdInFD)
	if err == nil {
		defer term.RestoreTerminal(stdInFD, oldInState)
	}

	oldOutState, err := term.SetRawTerminalOutput(stdOutFD)
	if err == nil {
		defer term.RestoreTerminal(stdOutFD, oldOutState)
	}

	// buff := &bytes.Buffer{}

	g := guy{}
	io.Copy(g, stdin)

	// for {
	// 	b, err := buff.ReadByte()
	// 	if err == nil {
	// 		fmt.Printf("%#v\n\r",b)
	// 	}
	// }	
}

type guy struct {
}

func (g guy) Write(p []byte) (n int, err error) {
	fmt.Printf("%#v\n\r", p)
	return len(p), nil
}