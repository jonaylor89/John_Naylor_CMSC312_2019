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
	Mem      memory.Memory
	InMsg    chan *Process
	ReadyQ   []*Process
	WaitingQ []*Process
	MinimumFreeFrames int
	// DeviceQ  []*Process
}

// RunRoundRobin : Start the schedule and process execution
func (s *Scheduler) RunRoundRobin() {

	// TimeQuantum : time quantum for process
	TimeQuantum := 50

	for {

		// Check for new processes
		s.recvProc()

		// Loop through processand execute them off a time quantum
		for i, curProc := range s.ReadyQ {
			curProc.state = RUN

			timeNull := s.CPU.TotalCycles

			// Only get so many CPU cycles
			for s.CPU.TotalCycles - timeNull < TimeQuantum {


				// Give the process access to the CPU and Process Channel
				err := curProc.Execute(s.CPU, s.InMsg)
				if err != nil {
					curProc.state = EXIT
					s.ReadyQ = remove(s.ReadyQ, i)
					s.Mem.RemovePages(curProc.PID)
					break
				}

				if curProc.runtime <= 0 {
					curProc.state = EXIT
					s.ReadyQ = remove(s.ReadyQ, i)
					break
				} else {
					curProc.state = READY
				}
			}

		}

		// Check if waiting processes can be moved to ready
		s.assessWaiting()

	}
}

// RunFirstComeFirstServe : First come first serve algorithm
func (s *Scheduler) RunFirstComeFirstServe() {

	for {

		// Check for new processes
		s.recvProc()

		var curProc *Process
		curProc, s.ReadyQ = s.ReadyQ[0], s.ReadyQ[1:]

		curProc.state = RUN

		// Execute process until it terminates
		for {

			curProc.Execute(s.CPU, s.InMsg)

			if curProc.runtime <= 0 {
				curProc.state = EXIT
				break
			}
		}

		// Check if waiting processes can be moved to ready
		s.assessWaiting()

	}
}

func (s *Scheduler) assessWaiting() {
	for i, proc := range s.WaitingQ {	
		if s.memoryCheck() {
			s.WaitingQ = remove(s.WaitingQ, i)

			proc.state = READY

			s.ReadyQ = append(s.ReadyQ, proc)
		}
	}
}

func (s *Scheduler) memoryCheck() bool {
	if cap(s.Mem.PhysicalMemory) - len(s.Mem.PhysicalMemory) > s.MinimumFreeFrames {
		return true
	}

	return false
}

func (s *Scheduler) recvProc() {

	// Check for new processes to schedule
	select {
	case x, ok := <-s.InMsg:
		if ok {

			if s.memoryCheck() {
				// If memory available then set to READY
				x.state = READY

				// New process ready to be executed
				s.ReadyQ = append(s.ReadyQ, x)
			} else {
				// If memory not available then set to WAIT
				x.state = WAIT

				// New process waiting for memory
				s.WaitingQ = append(s.WaitingQ, x)
			}

		} else {
			// Channel is closed to execution must exit
			return
		}
	default:
		// No new processes
		break
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
