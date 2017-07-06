package main

import (
	"fmt"
	"github.com/lyondhill/vtclean"
)

func main() {

	data := vtclean.Clean("\x1b[2K\x1b[1G\x1b[1G       0/960\x1b[2K\x1b[4G\x1b[4G 0/17333\x1b[4G 128/17333", false)
	fmt.Printf("%q\n", data)
}