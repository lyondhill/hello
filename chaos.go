
package main

import (
	"os"
	"os/exec"
	"fmt"
)

type cmd struct {
  cmd string
  args []string	
}

func main() {
	cmds := []cmd{
		cmd{ "nanobox", []string{"run", "ls"}},
		cmd{ "nanobox", []string{"stop"}},
		cmd{ "nanobox", []string{"start"}},
		cmd{ "nanobox", []string{"destroy"}},
		cmd{ "nanobox", []string{"run", "ls"}},
	}

	for _, cmd := range cmds {
		err := runCmd(cmd.cmd, cmd.args)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func runCmd(cmd string, args []string) error {
	fmt.Printf("\nrunning command: %s %v\n", cmd, args)

	c := exec.Command(cmd, args...)
	c.Stdin  = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}