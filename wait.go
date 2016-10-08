package main


import "os/exec"
import "fmt"
import "time"

func main() {
	cmd := exec.Command("sleep", "10")
	cmd.Start()
	go func(cmd *exec.Cmd) {
		time.Sleep(5 * time.Second)
		cmd.Wait()
		fmt.Println("routine_waited")
	}(cmd)
	cmd.Wait()
	fmt.Println("waited")
}