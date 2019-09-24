package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	// "strconv"
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

var (

	// ProcNum : PID for the highest process
	ProcNum int = 0
)

// Process : Running set of code
type Process struct {
	PID     int
	state   int
	runtime int
	memory  int
}

// Scheduler : Module to schedule process to run
type Scheduler struct {
	inMsg     chan Process
	processes []Process
}

// CreateProc : create a new process correctly
func CreateProc(runtime int, mem int) Process {

	ProcNum++

	return Process{
		PID:     ProcNum,
		state:   CREATED,
		runtime: runtime,
		memory:  mem,
	}
}

// Run : Start the schedule and process execution
func (s *Scheduler) Run() {
	for {

		// Check for new processes to schedule
		select {
		case x, ok := <-s.inMsg:
			if ok {
				s.processes = append(s.processes, x)
			} else {
				// Channel is closed to execution must exit
				fmt.Println("[INFO] exiting")
				return
			}
		}

		for _, curProc := range s.processes {
			curProc.state = RUNNING

			// TODO: Sleep for now instead of actual execution
			time.Sleep(200 * time.Millisecond)

			curProc.state = WAITING
		}

		time.Sleep(2000 * time.Millisecond)
	}
}

func main() {

	// Message channel between main kernel and scheduler
	ch := make(chan Process, 10)
	defer close(ch)

	s := Scheduler{
		inMsg:     ch,
		processes: []Process{},
	}

	// Run the scheduler
	go s.Run()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("OS Shell")
	fmt.Println("---------------------")

	for {

		fmt.Print("==> ")
		text, err := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		if err != nil {
			fmt.Println("failed to read user input")
		}

		switch text {
		case "new":
			p := CreateProc(rand.Intn(500)+1, rand.Intn(100)+1)
      ch <- p
      fmt.Println("processes: ", len(s.processes), "; queue: ", len(ch))
    case "len":
      fmt.Println("processes: ", len(s.processes), "; queue: ", len(ch))
		case "exit":
			fmt.Println("exiting simulator")
			return
    }
	}
}
