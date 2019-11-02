package sched

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (

	// TimeQuantum : time quantum for process
	TimeQuantum = 50

	// Process States

	// NEW : process created
	NEW = iota

	// READY : process in memory and ready for CPU
	READY

	// RUN : process running
	RUN

	// WAIT : process blocked
	WAIT

	// EXIT : process terminated
	EXIT
)

var (

	// ProcNum : PID for the highest process
	ProcNum int = 0
)

// Process : Running set of code
type Process struct {
	// Some info should be in a process contol block
	// And there will be a list of all process control blocks
	PID     int    // Process ID
	Name    string // Process Name
	state   int    // Process State
	runtime int    // Runtime Requirement
	memory  int    // Memory Requirement
}

// Scheduler : Controller to schedule process to run
type Scheduler struct {
	InMsg    chan *Process
	ReadyQ   []*Process
	WaitingQ []*Process
	// DeviceQ  []*Process
}

// CreateProcess : create a new process correctly
func CreateProcess(name string, runtime int, mem int) *Process {

	ProcNum++

	return &Process{
		PID:     ProcNum,
		Name:    name,
		state:   NEW,
		runtime: runtime,
		memory:  mem,
	}
}

// PickVictim : Pick a victum process to remove from physical memory
func PickVictim() {

}

// RunRoundRobin : Start the schedule and process execution
func (s *Scheduler) RunRoundRobin() {
	for {

		// Check for new processes to schedule
		select {
		case x, ok := <-s.InMsg:
			if ok {
				// New process ready to be executed
				s.ReadyQ = append(s.ReadyQ, x)

			} else {
				// Channel is closed to execution must exit
				return
			}
		default:
			// No new processes
			break
		}

		for i, curProc := range s.ReadyQ {
			curProc.state = RUN

			// I'm assuming this will get much more complex beyond just subtracting runtime
			// Fortunately, as of now it is basic round robin execution
			curProc.runtime -= TimeQuantum
			time.Sleep(100 * time.Millisecond)

			if curProc.runtime <= 0 {
				s.ReadyQ = remove(s.ReadyQ, i)
			} else {
				curProc.state = READY
			}

		}

	}
}

// CreateRandomProcessFromTemplate : Jitter template values to create custom processes
func CreateRandomProcessFromTemplate(templateName string, instructions [][]string, ch chan *Process) {

	r := rand.New(rand.NewSource(time.Now().Unix()))

	totalRuntime := 0
	for _, instruction := range instructions {
		if len(instruction) < 2 {
			continue
		}

		templateRuntime, err := strconv.Atoi(instruction[1])
		if err != nil {
			fmt.Println("error converting runtime to int", err)
		}

		// Jitter values by +-20
		templateRuntime += rand.Intn(20) - 10

		if templateRuntime < 0 {
			templateRuntime = 0
		}

		if instruction[0] == "CALCULATE" {
			totalRuntime += templateRuntime
		}

		instruction[1] = strconv.Itoa(templateRuntime)
	}

	p := CreateProcess("From template: "+templateName, totalRuntime, r.Intn(100)+1)
	p.state = READY
	ch <- p
}

func remove(slice []*Process, s int) []*Process {
	slice[s] = slice[len(slice)-1] // Copy last element to index i.
	// slice[len(slice)-1] = nil   	// Erase last element (write zero value)
	slice = slice[:len(slice)-1] // Truncate slice.

	return slice
}
