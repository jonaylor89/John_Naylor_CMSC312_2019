package sched

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (

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
	ch <- p
}
