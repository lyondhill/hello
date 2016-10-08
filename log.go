package main


import "fmt"
import "os/user"

func main() {
	username := "postgres"
	usr, _ := user.Current()
	if usr != nil {
		username = usr.Username
	}
	fmt.Println(username)
}