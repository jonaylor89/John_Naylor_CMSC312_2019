package kernel

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/cpu"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/memory"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/utils"
)

// Kernel : manager for resources and controller to schedule process to run
type Kernel struct {
	CPU               *cpu.CPU
	Mem               *memory.Memory
	InMsg             chan *Process
	ReadyQ            []*Process
	WaitingQ          []*Process
	MinimumFreeFrames int
	TimeQuantum int
	Mailboxes         []chan byte
	// DeviceQ  []*Process
}

// RunRoundRobin : Start the schedule and process execution
func (k *Kernel) RunRoundRobin() {

	// Check for new processes
	go k.recvProc()

	for {

		// Loop through processes backwards and execute them off a time quantum
		for i := len(k.ReadyQ) - 1; i > 0; i-- {
			curProc := k.ReadyQ[i]

			curProc.state = RUN

			timeNull := k.CPU.TotalCycles

			// Only get so many CPU cycles
			for k.CPU.TotalCycles-timeNull < k.TimeQuantum && !curProc.Critical {

				// Give the process access to the CPU and Process Channel
				err := curProc.Execute(k.CPU, k.Mem, k.InMsg, k.Mailboxes)
				if err != nil {

					curProc.state = EXIT
					k.ReadyQ = remove(k.ReadyQ, i)
					k.Mem.RemovePages(curProc.PID)
					break
				}

				curProc.state = READY
			}

		}

		// Check if waiting processes can be moved to ready
		k.assessWaiting()

	}
}

// RunFirstComeFirstServe : First come first serve algorithm
func (k *Kernel) RunFirstComeFirstServe() {

	// Check for new processes
	go k.recvProc()

	for {

		if len(k.ReadyQ) > 0 {

			var curProc *Process
			curProc, k.ReadyQ = k.ReadyQ[0], k.ReadyQ[1:]

			curProc.state = RUN

			// Execute process until it terminates
			for {

				curProc.Execute(k.CPU, k.Mem, k.InMsg, k.Mailboxes)

				if curProc.runtime <= 0 {
					curProc.state = EXIT
					break
				}
			}

		}

		// Check if waiting processes can be moved to ready
		k.assessWaiting()

	}
}

func (k *Kernel) assessWaiting() {
	for i := len(k.WaitingQ) - 1; i > 0; i-- {
		proc := k.WaitingQ[i]

		if k.memoryCheck() {
			k.WaitingQ = remove(k.WaitingQ, i)

			proc.state = READY

			k.ReadyQ = append(k.ReadyQ, proc)

			k.Mem.Add(proc.memory, proc.PID)
		}
	}
}

func (k *Kernel) memoryCheck() bool {
	if cap(k.Mem.PhysicalMemory)-len(k.Mem.PhysicalMemory) > k.MinimumFreeFrames {
		return true
	}

	return false
}

func (k *Kernel) recvProc() {

	for {

		// Check for new processes to schedule
		select {
		case x, ok := <-k.InMsg:
			if ok {

				if k.memoryCheck() {

					// If memory available then set to READY
					x.state = READY

					// New process ready to be executed
					k.ReadyQ = append(k.ReadyQ, x)

				} else {
					// If memory not available then set to WAIT
					x.state = WAIT

					// New process waiting for memory
					k.WaitingQ = append(k.WaitingQ, x)

				}

				x.pages = k.Mem.Add(x.memory, x.PID)

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
