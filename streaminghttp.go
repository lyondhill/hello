package main

import (
	"net"
	"net/http"
	// "net/http/httputil"
	"os"
	"io"
	"time"
	"fmt"
	"runtime"
	gosignal "os/signal"

	"github.com/docker/docker/pkg/term"
	"github.com/docker/docker/pkg/signal"
)

func main() {
	conn, err := net.Dial("tcp4", "nanobox-ruby-sample.nano.dev:1757")
	if err != nil {
		fmt.Println(err)
		return
	}
	// conn.SetDeadline(time.Now())
	// forward all the signals to the nanobox server
	forwardAllSignals()

	// make sure we dont just kill the connection
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}
	// fake a web request
	conn.Write([]byte("POST /exec HTTP/1.1\r\n\r\n"))

	// setup a raw terminal
	var oldState *term.State
	stdIn, stdOut, _ := term.StdStreams()
	inFd, _ := term.GetFdInfo(stdIn)
	outFd, _ := term.GetFdInfo(stdOut)
	
	// monitor the window size and send a request whenever we resize
	monitorSize(outFd)

	oldState, err = term.SetRawTerminal(inFd)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer term.RestoreTerminal(inFd, oldState)

	// pipe data
	go io.Copy(os.Stdout, conn)
	io.Copy(conn, os.Stdin)
}

// forward the signals you recieve and send them to nanobox server
func forwardAllSignals() {
	sigc := make(chan os.Signal, 128)
	// signal.CatchAll(sigc)
	go func() {
		for s := range sigc {
			// skip children and window resizes
			continue
			if s == signal.SIGCHLD || s == signal.SIGWINCH {
				
			}
			var sig string
			for sigStr, sigN := range signal.SignalMap {
				if sigN == s {
					sig = sigStr
					break
				}
			}
			if sig == "" {
				fmt.Printf("Unsupported signal: %v. Discarding.\n", s)
			}
			fmt.Println(sig)
			req, _ := http.NewRequest("POST", fmt.Sprintf("http://nanobox-ruby-sample.nano.dev:1757/killexec?signal=%s", sig), nil)
			_, err := http.DefaultClient.Do(req)
			fmt.Println(err)
		}
	}()
	return
}

func monitorSize(outFd uintptr) {
	resizeTty(outFd)

	if runtime.GOOS == "windows" {
		go func() {
			prevH, prevW := getTtySize(outFd)
			for {
				time.Sleep(time.Millisecond * 250)
				h, w := getTtySize(outFd)

				if prevW != w || prevH != h {
					resizeTty(outFd)
				}
				prevH = h
				prevW = w
			}
		}()
	} else {
		sigchan := make(chan os.Signal, 1)
		gosignal.Notify(sigchan, signal.SIGWINCH)
		go func() {
			for range sigchan {
				resizeTty(outFd)
			}
		}()
	}	
}

func resizeTty(outFd uintptr) {
	h, w := getTtySize(outFd)
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://nanobox-ruby-sample.nano.dev:1757/resizeexec?h=%d&w=%d", h, w), nil)
	http.DefaultClient.Do(req)
}

func getTtySize(outFd uintptr) (h, w int) {
	ws, err := term.GetWinsize(outFd)
	if err != nil {
		fmt.Printf("Error getting size: %s\n", err)
		if ws == nil {
			return 0, 0
		}
	}
	h = int(ws.Height)
	w = int(ws.Width)
	return
}
