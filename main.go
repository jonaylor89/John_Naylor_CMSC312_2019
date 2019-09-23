
package main

import (
  "fmt"
  "time"
)

const (

  // ClockSpeed : Execution Rate for virtual CPU
  ClockSpeed = 200

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

var (

  // NumOfProc : Number of processes created
  NumOfProc int = 0
)

// Process : Running set of code
type Process struct {
  PID int;
  state int;
  runtime int;
  memory int;
}

func createProc(runtime int, mem int) Process {

  NumOfProc++

  return Process {
    PID: NumOfProc,
    state: CREATED,
    runtime: runtime,
    memory: mem,
  }
}

func main() {

  p := createProc(300, 45)

  for {

    p.state = RUNNING
    fmt.Println(p)
    p.state = WAITING

    time.Sleep(ClockSpeed * time.Millisecond)

  }

}
