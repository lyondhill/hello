// The primary mechanism for managing state in Go is
// communication over channels. We saw this for example
// with [worker pools](worker-pools). There are a few other
// options for managing state though. Here we'll
// look at using the `sync/atomic` package for _atomic
// counters_ accessed by multiple goroutines.

package main

import "fmt"
import "time"
import "sync/atomic"
import "runtime"

func main() {

    // We'll use an unsigned integer to represent our
    // (always-positive) counter.
    var ops uint32 = 0
    ops = ops -1
    // To simulate concurrent updates, we'll start 50
    // goroutines that each increment the counter about
    // once a millisecond.
    for i := 0; i < 50; i++ {
        go func() {
            for {
                // To atomically increment the counter we
                // use `AddUint64`, giving it the memory
                // address of our `ops` counter with the
                // `&` syntax.
                atomic.AddUint32(&ops, 1)

                // Allow other goroutines to proceed.
                runtime.Gosched()
            }
        }()
    }

    // Wait a second to allow some ops to accumulate.
    time.Sleep(time.Second)

    opsFinal := atomic.LoadUint32(&ops)
    fmt.Println("ops:", opsFinal)
    time.Sleep(time.Second)

    opsFinal = atomic.LoadUint32(&ops)
    fmt.Println("ops:", opsFinal)
    time.Sleep(time.Second)

    opsFinal = atomic.LoadUint32(&ops)
    fmt.Println("ops:", opsFinal)
    time.Sleep(time.Second)

    opsFinal = atomic.LoadUint32(&ops)
    fmt.Println("ops:", opsFinal)
    time.Sleep(time.Second)

    opsFinal = atomic.LoadUint32(&ops)
    fmt.Println("ops:", opsFinal)
}
