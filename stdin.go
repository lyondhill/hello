package main

import "github.com/nsf/termbox-go"
import "fmt"

// import "bytes"
// import "os/exec"
// import "os"
// import "io"

func main() {



	// buffer := &bytes.Buffer{}
	// go io.Copy(buffer, os.Stdin)
	// for {
	// 	b := make([]byte, 10)
	// 	n, _ := buffer.Read(b)
	// 	if n == 0 {
	// 		continue
	// 	}
	// 	fmt.Printf("%#v\n\n", b)
	// }

	// cmd := exec.Command("winpty", "--showkey")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// fmt.Println(cmd.Start())
	// cmdin, _ := cmd.StdinPipe()
	// io.Copy(cmdin, os.Stdin)
	// cmd.Wait()

	fmt.Println()
	termbox.Init()
	fmt.Println([]byte("^[[A"))

	for {
		event := termbox.PollEvent()
		if event.Key == termbox.KeyCtrlC {
			termbox.Close()
			return
		}
		fmt.Printf("%#v\n", event)		
	}
}