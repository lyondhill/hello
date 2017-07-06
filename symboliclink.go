package main

import "fmt"
import "path/filepath"
import "os"

func main() {
	dir, _ := os.Getwd()
	fmt.Println(filepath.EvalSymlinks(dir))
}