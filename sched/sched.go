package sched

import (
	"bufio"
	// "fmt"
	"io"
	"os"
	"strconv"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/cpu"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/memory"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/utils"
)

// Scheduler : manager for resources and controller to schedule process to run
type Scheduler struct {
	CPU               *cpu.CPU
	Mem               *memory.Memory
	InMsg             chan *Process
	ReadyQ            []*Process
	WaitingQ          []*Process
	MinimumFreeFrames int
	TimeQuantum       int
	Mailboxes         []chan byte
	// DeviceQ  []*Process
}

// InitScheduler : create new scheduler
func InitScheduler(cpu *cpu.CPU, mem *memory.Memory, in chan *Process,minimumFreeFrames int, timeQuantum int) *Scheduler {
	return &Scheduler{
		CPU:               cpu,
		Mem:               mem,
		InMsg:             in,
		ReadyQ:            []*Process{},
		WaitingQ:          []*Process{},
		MinimumFreeFrames: minimumFreeFrames,
		TimeQuantum:       timeQuantum,
		Mailboxes: []chan byte{
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
		},
	}
}

// RunRoundRobin : Start the schedule and process execution
func (s *Scheduler) RunRoundRobin() {

	// Check for new processes
	go s.recvProc()

	for {

		// Loop through processes backwards and execute them off a time quantum
		for i := len(s.ReadyQ) - 1; i > 0; i-- {
			curProc := s.ReadyQ[i]

			curProc.State = RUN

			timeNull := s.CPU.TotalCycles

			// Only get so many CPU cycles
			for s.CPU.TotalCycles-timeNull < s.TimeQuantum && !curProc.Critical {

				// Give the process access to the CPU and Process Channel
				err := curProc.Execute(s.CPU, s.Mem, s.InMsg, s.Mailboxes)
				if err != nil {

					curProc.State = EXIT
					s.ReadyQ = remove(s.ReadyQ, i)
					s.Mem.RemovePages(curProc.PID)
					break
				}

				curProc.State = READY
			}

		}

		// Check if waiting processes can be moved to ready
		s.assessWaiting()

	}
}

// RunFirstComeFirstServe : First come first serve algorithm
func (s *Scheduler) RunFirstComeFirstServe() {

	// Check for new processes
	go s.recvProc()

	for {

		if len(s.ReadyQ) > 0 {

			var curProc *Process
			curProc, s.ReadyQ = s.ReadyQ[0], s.ReadyQ[1:]

			curProc.State = RUN

			// Execute process until it terminates
			for {

				curProc.Execute(s.CPU, s.Mem, s.InMsg, s.Mailboxes)

				if curProc.Runtime <= 0 {
					curProc.State = EXIT
					break
				}
			}

		}

		// Check if waiting processes can be moved to ready
		s.assessWaiting()

	}
}

func (s *Scheduler) assessWaiting() {
	for i := len(s.WaitingQ) - 1; i > 0; i-- {
		proc := s.WaitingQ[i]

		if s.memoryCheck() {
			s.WaitingQ = remove(s.WaitingQ, i)

			proc.State = READY

			s.ReadyQ = append(s.ReadyQ, proc)

			s.Mem.Add(proc.Memory, proc.PID)
		}
	}
}

func (s *Scheduler) memoryCheck() bool {
	if cap(s.Mem.PhysicalMemory)-len(s.Mem.PhysicalMemory) > s.MinimumFreeFrames {
		return true
	}

	return false
}

func (s *Scheduler) recvProc() {

	for {

		// Checs for new processes to schedule
		select {
		case x, ok := <-s.InMsg:
			if ok {

				if s.memoryCheck() {

					// If memory available then set to READY
					x.State = READY

					// New process ready to be executed
					s.ReadyQ = append(s.ReadyQ, x)

				} else {
					// If memory not available then set to WAIT
					x.State = WAIT

					// New process waiting for memory
					s.WaitingQ = append(s.WaitingQ, x)

				}

				x.pages = s.Mem.Add(x.Memory, x.PID)

			} else {
				// Channel is closed to execution must exit
				return
			}
		default:
			// No new processes
			break
		}

	}
}

// LoadTemplate : load in template process and create process mutations off of it
func LoadTemplate(filename string, numOfProcesses int, processChan chan *Process) error {

	f, err := os.Open(filename)
	if err != nil {
		// fmt.Println("error opening file", err)
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
				// fmt.Printf(" > Failed!: %v\n", err)
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
