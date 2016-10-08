package main


import "syscall"
import "fmt"

func main() {
	rlm := syscall.Rlimit{}
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlm)
	fmt.Printf("cur: %d, max: %d",rlm.Cur, rlm.Max)
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlm)
}