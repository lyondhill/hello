package main

import (
	"context"
	"fmt"
	"time"

	"github.com/influxdata/platform/query"
	_ "github.com/influxdata/platform/query/builtin"
)

func main() {
	script := `option task = {
      name: "name",
      concurrency: 1,
      every: 1m0s,
}

from(db: "test")
    |> range(start:-1h)`
	inter := query.NewInterpreter()
	fmt.Println("eval", query.Eval(inter, script))

	fmt.Println("inter", inter)

	spec, err := query.Compile(context.Background(), script, time.Now())
	spec.Walk(func(o *query.Operation) error {
		fmt.Printf("op: %#v\n", o)
		return nil
	})
	fmt.Printf("spec: %+v\n", spec)
	fmt.Println(err)

}
