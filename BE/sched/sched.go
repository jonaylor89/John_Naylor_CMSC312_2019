package sched

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"os"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/memory"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/utils"
)

// Scheduler : Controller to schedule process to run
type Scheduler struct {
	CPU      CPU
	RAM      memory.RAM
	InMsg    chan *Process
	ReadyQ   []*Process
	WaitingQ []*Process
	// DeviceQ  []*Process
}

// RunRoundRobin : Start the schedule and process execution
func (s *Scheduler) RunRoundRobin() {

	// TimeQuantum : time quantum for process
	TimeQuantum := 50

	for {

		// Check for new processes to schedule
		select {
		case x, ok := <-s.InMsg:
			if ok {

				x.state = READY

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
			for j := 0; j < TimeQuantum; j++ {


				curProc.Execute(s.CPU, s.InMsg)

				if curProc.runtime <= 0 {
					s.ReadyQ = remove(s.ReadyQ, i)
					break
				} else {
					curProc.state = READY
				}
			}

		}

	}
}

// RunFirstComeFirstServe : First come first serve algorithm
func (s *Scheduler) RunFirstComeFirstServe() {

	for {

		// Check for new processes to schedule
		select {
		case x, ok := <-s.InMsg:
			if ok {

				x.state = READY

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

		var curProc *Process
		curProc, s.ReadyQ = s.ReadyQ[0], s.ReadyQ[1:]

		curProc.state = RUN

		for {

			curProc.Execute(s.CPU, s.InMsg)

			if curProc.runtime <= 0 {
				break
			}
		}

	}
}

// LoadTemplate : load in template process and create process mutations off of it
func LoadTemplate(filename string, numOfProcesses int, processChan chan *Process) error {

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file", err)
		return nil
	}

	defer f.Close()

	reader := bufio.NewReader(f)

	var instructions [][]string

	// Get process name from the first line
	procNameField, _ := utils.ReadLine(reader)
	procName := procNameField[1]

	// Get the process memory requirement from the second line
	procMemoryField, _ := utils.ReadLine(reader)
	procMemory, _ := strconv.Atoi(procMemoryField[1])

	// Loop through template file from the instructions
	for {

		instruction, err := utils.ReadLine(reader)
		if err != nil {
			if err != io.EOF {
				fmt.Printf(" > Failed!: %v\n", err)
				return err
			}

			break
		}

		if len(instruction) != 0 {

			instructions = append(instructions, instruction)
		}
	}

	// Randomize order of isntructions
	// utils.ShuffleInstructions(instructions)
	
	for i := 0; i < numOfProcesses; i++ {
		go CreateRandomProcessFromTemplate(procName, procMemory, instructions, processChan)
	}

	return nil
}

func remove(slice []*Process, s int) []*Process {
	slice[s] = slice[len(slice)-1] // Copy last element to index i.
	// slice[len(slice)-1] = nil   	// Erase last element (write zero value)
	slice = slice[:len(slice)-1] // Truncate slice.

	return slice
}
