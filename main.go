package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
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

	// TimeQuantum : time quantum for process
	TimeQuantum = 10

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

// Scheduler : Module to schedule process to run
type Scheduler struct {
	inMsg     chan *Process
	processes []*Process
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

func remove(slice []*Process, s int) []*Process {
	return append(slice[:s], slice[s+1:]...)
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

func main() {

	rand.Seed(time.Now().UnixNano())

	// Message channel between main kernel and scheduler
	ch := make(chan *Process, 10)
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

		case "new":
			p := CreateProc("Random Proc", rand.Intn(500)+1, rand.Intn(100)+1)
			ch <- p
			fmt.Println("processes: ", len(s.processes), "; queue: ", len(ch))

		case "load":
			if len(args) != 2 {
				fmt.Println("`load` requires a filename as an argument")
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
			for {
				line, err = reader.ReadString('\n')
				if err != nil {
					break
				}

				fmt.Printf(" > Read %d characters\n", len(line))
				fmt.Println(line)
			}

			if err != io.EOF {
				fmt.Printf(" > Failed!: %v\n", err)
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
