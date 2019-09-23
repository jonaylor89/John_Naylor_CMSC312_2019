
package main

import (
  "fmt"
  "time"
)

const (

  // Process States

  // CREATED : process created
  CREATED = iota

  // RUNNING : process running
  RUNNING

  // WAITING : process waiting
  WAITING

  // BLOCKED : process blocked
  BLOCKED

  // TERMINATED : process terminated
  TERMINATED
)

// Process : Running set of code
type Process struct {
  PID int;
  state int;
  runtime int;
  memory int;
}

func main() {

  p := Process{ 
    PID: 1, 
    state: CREATED,
    runtime: 300,
    memory: 45,
   }

  for {

    p.state = RUNNING
    fmt.Println(p)
    p.state = WAITING

    time.Sleep(200 * time.Millisecond)

  }

}
