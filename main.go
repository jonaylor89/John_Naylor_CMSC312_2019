package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	// "sort"
	"strconv"
	"strings"
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

	// TimeQuantum : time quantum for process
	TimeQuantum = 50

	// ProcNum : PID for the highest process
	ProcNum int = 0
)

// Process : Running set of code
type Process struct {
	PID     int
	Name    string
	state   int
	runtime int
	memory  int
}

// Scheduler : Controller to schedule process to run
type Scheduler struct {
	inMsg     chan *Process
	processes []*Process
}

// Run : Start the schedule and process execution
func (s *Scheduler) Run() {
	for {

		// Check for new processes to schedule
		select {
		case x, ok := <-s.inMsg:
			if ok {
				// New process ready to be executed
				s.processes = append(s.processes, x)

			} else {
				// Channel is closed to execution must exit
				return
			}
		default:
			// No new processes
			break
		}

		for i, curProc := range s.processes {
			curProc.state = RUNNING

			// I'm assuming this will get much more complex beyond just subtracting runtime
			// Fortunately, as of now it is basic round robin execution
			curProc.runtime -= TimeQuantum
			time.Sleep(200 * time.Millisecond)

			if curProc.runtime <= 0 {
				s.processes = remove(s.processes, i)
			} else {
				curProc.state = WAITING
			}

		}

	}
}


// CreateProc : create a new process correctly
func CreateProc(name string, runtime int, mem int) *Process {

	ProcNum++

	return &Process{
		PID:     ProcNum,
		Name:    name,
		state:   CREATED,
		runtime: runtime,
		memory:  mem,
	}
}

func createRandomProcessFromTemplate(templateName string, instructions [][]string, ch chan *Process) {

	totalRuntime := 0
	for _, instruction := range instructions {
		if len(instruction) < 2 {
			continue
		}



		templateRuntime, err := strconv.Atoi(instruction[1])
		if err != nil {
			fmt.Println("Error converting runtime to int", err)
		}

		// Jitter values by +-10
		templateRuntime += rand.Intn(20) - 10

		if instruction[0] == "CALCULATE" {
			totalRuntime += templateRuntime
		}

		instruction[1] = strconv.Itoa(templateRuntime)
	}

	p := CreateProc("From template: "+templateName, totalRuntime, rand.Intn(100)+1)
	ch <- p
}

func shuffleInstructions(vals [][]string) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	// We start at the end of the slice, inserting our random
	// values one at a time.
	for n := len(vals); n > 0; n-- {
		randIndex := r.Intn(n)
		// We swap the value at index n-1 and the random index
		// to move our randomly chosen value to the end of the
		// slice, and to move the value that was at n-1 into our
		// unshuffled portion of the slice.
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
	}
}

func remove(slice []*Process, s int) []*Process {
	slice[s] = slice[len(slice)-1]  // Copy last element to index i.
	// slice[len(slice)-1] = nil   	// Erase last element (write zero value)
	slice = slice[:len(slice)-1]   	// Truncate slice.

	return slice
}

func main() {

	rand.Seed(time.Now().UnixNano())

	// Message channel between main kernel and scheduler
	ch := make(chan *Process, 1000)
	defer close(ch)

	s := Scheduler{
		inMsg:     ch,
		processes: []*Process{},
	}

	// Run the scheduler
	go s.Run()

	console := bufio.NewReader(os.Stdin)
	fmt.Println("OS Shell")
	fmt.Println("---------------------")

	var args []string
	for {

		fmt.Print("==> ")
		text, err := console.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		args = strings.Split(text, " ")
		if err != nil {
			fmt.Println("failed to read user input")
		}

		switch args[0] {

		case "load":
			if len(args) != 3 {
				fmt.Println("`load` requires a filename and number of processes as an argument")
				break
			}

			f, err := os.Open(args[1])
			if err != nil {
				fmt.Println("error loading file")
				break
			}

			defer f.Close()

			reader := bufio.NewReader(f)

			var line string
			var instruction []string
			var instructions [][]string

			for {
				line, err = reader.ReadString('\n')
				if err != nil {
					break
				}

				line = strings.ReplaceAll(line, "\n", "")
				instruction = strings.Split(line, " ")

				if len(instruction) != 2 || (instruction[0] != "CALCULATE" && instruction[0] != "I/O") {
					// Skip the first few lines with meta data and only work with instructions for now
					continue
				}

				instructions = append(instructions, instruction)

			}

			if err != io.EOF {
				fmt.Printf(" > Failed!: %v\n", err)
				break
			}

			numOfProc, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println("Could not get number of processes")
				break
			}

			// Randomize order of isntructions
			shuffleInstructions(instructions)

			for i := 0; i < numOfProc; i++ {
				go createRandomProcessFromTemplate(args[1], instructions, ch)
			}

		case "len":
			fmt.Println("processes: ", len(s.processes), "; queue: ", len(ch))

		case "dump":
			fmt.Println("process dump:")
			for _, proc := range s.processes {
				fmt.Println(*proc)
			}

		case "exit":
			fmt.Println("exiting simulator")
			return

		default:
			break

		}
	}
}
