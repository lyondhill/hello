package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	// read file
	data, err := ioutil.ReadFile(os.Getenv("HOME") + "/.pcl.toml")
	if err != nil {
		panic(err)
	}

	s := struct {
		Owner      string `toml:"owner"`
		Key        string `toml:"key"`
		LicenceKey string `toml:"enterprise-license-key"`
		AWS        struct {
			Key    string `toml:"access-key-id"`
			Secret string `toml:"secret-access-key"`
			Region string `toml:"region"`
		} `toml:"aws"`
	}{}

	toml.Unmarshal(data, &s)
	fmt.Printf("%+v\n", s)

	homeDir := os.Getenv("HOME")
	keyPath := strings.Replace(s.Key, "${HOME}", homeDir, 1)
	keyPath = strings.Replace(keyPath, "~", homeDir, 1)
	data, err = ioutil.ReadFile(keyPath)
	fmt.Printf("%s, %v\n\n", data, err)
	encoder := toml.NewEncoder(os.Stdout)
	encoder.Encode(s)
}
