
package main

import (
  "fmt"
  "time"
  "math/rand"
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

  // ProcNum : PID for the highest process
  ProcNum int = 0
)

// Process : Running set of code
type Process struct {
  PID int;
  state int;
  runtime int;
  memory int;
}

// Scheduler : Module to schedule process to run
type Scheduler struct {
  processes []Process;
}

func createProc(runtime int, mem int) Process {

  ProcNum++

  return Process {
    PID: ProcNum,
    state: CREATED,
    runtime: runtime,
    memory: mem,
  }
}

func main() {

  s := Scheduler{
    processes: []Process{},
  }

  p := createProc(rand.Intn(500) + 1, rand.Intn(100) + 1)

  s.processes = append(s.processes, p)

  for {

    for _, curProc := range s.processes {
      curProc.state = RUNNING

      fmt.Println(curProc)


      curProc.state = WAITING
      time.Sleep(ClockSpeed * time.Millisecond)
    }

  }

}
