package main

import "os"
import "github.com/fsnotify/fsnotify"
import "path/filepath"
import "log"

var watcher *fsnotify.Watcher

func main() {

  var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
	    log.Fatal(err)
	}

  filepath.Walk(".", fileWalker)

  for {
    
    select {
    case event := <-watcher.Events:

      watcher.Add(event.Name)
      log.Println("event:", event)
      if event.Op&fsnotify.Write == fsnotify.Write {
        log.Println("modified file:", event.Name)
      }
    case err := <-watcher.Errors:
      log.Println("error:", err)
    }
  }
	
}

func fileWalker(path string, info os.FileInfo, err error) error {
  if err != nil {
    return err
  }
  err = watcher.Add(path)
  return err
}

