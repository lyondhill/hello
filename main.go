package main

import (
    "code.google.com/p/gcfg"
    "fmt"
    "log"
)

func main() {
    cfgStr := `; Comment line
[section]
var-name=value # comment`
    cfg := struct {
        Section struct {
            FieldName string `gcfg:"var-name"`
        }
    }{}
    err := gcfg.ReadStringInto(&cfg, cfgStr)
    if err != nil {
        log.Fatalf("Failed to parse gcfg data: %s", err)
    }
    fmt.Println(cfg.Section.FieldName)
}
