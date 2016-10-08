package main

import (
	"fmt"
	"os"
	"path/filepath"
)


func main() {
	fmt.Println("oldLocalDir", oldLocalDir())
	fmt.Println("localDir   ", LocalDir())
}
// oldLocalDir ...
func oldLocalDir() string {

	//
	p, err := os.Getwd()
	if err != nil {
		return ""
	}

	return filepath.ToSlash(p)
}

// LocalDir ...
func LocalDir() string {

	//
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	boxfilePresent := func(path string) bool {
		boxfile := filepath.ToSlash(filepath.Join(path, "boxfile.yml"))
		fi, err := os.Stat(boxfile)
		if err != nil {
			return false
		}
		return !fi.IsDir()
	}

	path := cwd
	for !boxfilePresent(path) {
		if path == "" || path == "/" {
			// return the current working directory if we cant find a path
			return cwd
		}
		// eliminate the most child directory and then check it
		path = filepath.Dir(path)
	}

	// recursively check for boxfile


	return filepath.ToSlash(path)
}