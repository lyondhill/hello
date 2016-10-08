package main

import (
  "fmt"
  "os"
  "io"

  // "github.com/docker/docker/pkg/stdcopy"
  "github.com/docker/engine-api/client"
  "github.com/docker/engine-api/types"
  "golang.org/x/net/context"
  "github.com/docker/docker/pkg/term"

)

func main() {
  client, err := client.NewEnvClient()
  fmt.Println(err)

  fmt.Println(client.ClientVersion())
  fmt.Println(client.ServerVersion(context.Background()))
  config := types.ExecConfig{
    Tty:          true,
    Cmd:          []string{"/bin/bash"},
    AttachStdin:  true,
    AttachStdout: true,
    AttachStderr: true,
  }

  exec, err := client.ContainerExecCreate(context.Background(), "8bbe872bac96", config)
  fmt.Println(err)

  resp, err := client.ContainerExecAttach(context.Background(), exec.ID, config)
  fmt.Println(err)

  stdIn, stdOut, _ := term.StdStreams()
  stdInFD, _ := term.GetFdInfo(stdIn)
  stdOutFD, _ := term.GetFdInfo(stdOut)

  oldInState, err := term.SetRawTerminal(stdInFD)
  oldOutState, err := term.SetRawTerminalOutput(stdOutFD)

  defer term.RestoreTerminal(stdInFD, oldInState)
  defer term.RestoreTerminal(stdOutFD, oldOutState)


  // fmt.Println(ExecPipe(resp, os.Stdin, os.Stdout, os.Stderr))
  go io.Copy(resp.Conn, os.Stdin)
  io.Copy(os.Stdout, resp.Reader)
}

func ExecPipe(resp types.HijackedResponse, inStream io.Reader, outStream, errorStream io.Writer) error {
  var err error
  receiveStdout := make(chan error, 1)
  if outStream != nil || errorStream != nil {
    go func() {
      // always do this because we are never tty
      // _, err = stdcopy.StdCopy(outStream, errorStream, resp.Reader)
      _, err = io.Copy(outStream, resp.Reader)
      fmt.Printf("[hijack] End of stdout")
      receiveStdout <- err
    }()
  }

  stdinDone := make(chan struct{})
  go func() {
    if inStream != nil {
      io.Copy(resp.Conn, inStream)
      fmt.Printf("[hijack] End of stdin")
    }

    if err := resp.CloseWrite(); err != nil {
      fmt.Printf("Couldn't send EOF: %s", err)
    }
    close(stdinDone)
  }()

  select {
  case err := <-receiveStdout:
    if err != nil {
      fmt.Printf("Error receiveStdout: %s", err)
      return err
    }
  case <-stdinDone:
    if outStream != nil || errorStream != nil {
      if err := <-receiveStdout; err != nil {
        fmt.Printf("Error receiveStdout: %s", err)
        return err
      }
    }
  }

  return nil
}
