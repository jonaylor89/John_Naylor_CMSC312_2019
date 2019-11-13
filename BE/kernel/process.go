package kernel

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/code"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/memory"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/cpu"
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
	PID      int    // Process ID
	Name     string // Process Name
	state    int    // Process State
	runtime  int    // Runtime Requirement
	memory   int    // Memory Requirement
	children []int  // List of PID to child processes
	ip		 int    // Instruction pointer
	ins 	 code.Instructions 
	pages 	 []int  // memory pages owned by process
	Critical bool   // is the process in the critical section
}

// CreateProcess : create a new process correctly
func CreateProcess(name string, runtime int, mem int, ins code.Instructions, insPointer int) *Process {

	ProcNum++

	return &Process{
		PID:      ProcNum,
		Name:     name,
		state:    NEW,
		runtime:  runtime,
		memory:   mem,
		children: []int{},
		ip: 	  insPointer,
		ins: 	  ins,
		pages: 	  []int{},
		Critical: false,
	}
}

// String : string reporesentation of process
// func (p Process) String() string {
// 	return fmt.Sprintf("Name: %s ;Instructions: %s", 
// 					p.Name,
// 					p.ins)
// }

// Execute : execute instruction in process
func (p *Process) Execute(cpu *cpu.CPU, mem *memory.Memory, ch chan *Process) error {


	if len(p.ins) <= p.ip {
		return fmt.Errorf("End of instructions")
	}

	curIns := p.ins[p.ip]
	op := code.Opcode(curIns)

	switch (op) {

	case code.CALC:

		cpu.RunCycle(p.runtime)

		// Subtract one from the time
		p.ins[p.ip+1]--

		time := code.ReadUint8(p.ins[p.ip+1:])

		// Check if instruction is finished
		if time <= 0 {
			p.ip += 2
		}

		break
	case code.IO:
		p.ip += 2

		break
	case code.FORK:

		p.ip += 1

		// create child process
		child := CreateProcess("Fork: "+p.Name, p.runtime, p.memory, p.ins, p.ip)

		// Add child process to list of children of parent
		p.children = append(p.children, child.PID)

		// Send child to scheduler
		ch <- child

		break
	case code.ENTER:
		p.ip += 1

		p.Critical = true
		break
	case code.EXIT:
		p.ip += 1

		p.Critical = false
		break
	case code.NOP:
		p.ip += 1
		break
	default:
		p.ip += 1
		break
	}

	return nil
}

// CreateRandomProcessFromTemplate : Jitter template values to create custom processes
func CreateRandomProcessFromTemplate(templateName string, memory int, instructions [][]string, ch chan *Process) {

	r := rand.New(rand.NewSource(time.Now().Unix()))

	totalRuntime := 0
	for _, instruction := range instructions {
		if len(instruction) < 2 {
			continue
		}

		templateValue, err := strconv.Atoi(instruction[1])
		if err != nil {
			fmt.Println("error converting runtime to int", err)
		}

		// Jitter values by +-20
		templateValue += r.Intn(20) - 10

		if templateValue < 0 {
			templateValue = 0
		}

		if instruction[0] == "CALC" {
			totalRuntime += templateValue
		}

		instruction[1] = strconv.Itoa(templateValue)
	}

	program := code.Assemble(instructions)
	
	p := CreateProcess("From template: "+templateName, totalRuntime, memory, program, 0)
	ch <- p
}


