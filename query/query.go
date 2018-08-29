package main

import (
	"fmt"

	"github.com/influxdata/platform/query"
	_ "github.com/influxdata/platform/query/builtin"
	"github.com/influxdata/platform/query/parser"
	"github.com/influxdata/platform/query/semantic"
)

func main() {
	queryString := `option task = {
  // name is required, and must be unique for a give user_id & organization_id
  name: "foo",
  // every is a duration, this task should be run at this interval. minimum interval is 1s.
  every: 1h,
  // delay is a duration, in this example it would delay scheduling the task to
  // 10 minutes past the hour
  delay: 10m,
  // cron is a more sophisticated way to schedule. every and cron are mutually exclusive
  cron: "0 2 * * *", // run at 2 AM
  // retry is the number of times to retry a failed run of a task before
  // giving up. This isn't in scope of initial implementation
  retry: 5,
}

option now = ()
		=>
// Now in the rest of the script the values of task can be accessed.

from(db: "test")
    |> range(start:-task.every)
   //...
  `

	inter := query.NewInterpreter()
	fmt.Println(inter)

	ast, err := parser.NewAST(queryString)
	fmt.Println("ast", ast, err)

	_, declarations := query.BuiltIns()

	semanticProgram, err := semantic.New(ast, declarations)
	fmt.Println("sema", semanticProgram, err)

	// Evaluate program
	err = inter.Eval(semanticProgram)
	fmt.Println("evalerr", err)
	fmt.Println(inter)
	retry, _ := inter.Option("task").Object().Get("retry")
	retryval := retry.UInt()
	fmt.Printf("%+v\n", retryval)
}
