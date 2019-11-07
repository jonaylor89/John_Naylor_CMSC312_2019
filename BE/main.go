package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	// "time"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/sched"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/memory"
)

func main() {

	// Message channel between main kernel and scheduler
	ch := make(chan *sched.Process, 1000)
	defer close(ch)

	cpu := &sched.CPU{ 
		TotalCycles: 0, 
		Speed: 10,
	}

	mem := &memory.Memory{
		PageSize: 32,
		TotalRam: 4096,
		PageTable: make(map[int]int),
		VirtualMemory: make([]*memory.Page, 0),
		PhysicalMemory: make([]*memory.Page, 0, 4096 / 32),
	}

	s := &sched.Scheduler{
		CPU: 	  cpu,
		Mem: 	  mem,
		InMsg:    ch,
		ReadyQ:   []*sched.Process{},
		WaitingQ: []*sched.Process{},
		MinimumFreeFrames: 8,
	}

	// Run the scheduler
	go s.RunRoundRobin()
	// go s.RunFirstComeFirstServer()

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

			filename := args[1]

			numOfProc, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println("Could not get number of processes")
				break
			}

			if numOfProc <= 0 {
				fmt.Println("`load` number of processes must be postive")
				break
			}

			err = sched.LoadTemplate(filename, numOfProc, ch)
			if err != nil {
				fmt.Println("`load` error loading process template", err)
				break
			}

		case "proc":
			fmt.Println("ready: ", len(s.ReadyQ), "; waiting: ", len(s.WaitingQ), "; sending: ", len(ch))

		case "mem":
			fmt.Println("Physical len: ", len(s.Mem.PhysicalMemory), "; cap: ", cap(s.Mem.PhysicalMemory))
			fmt.Println("Virtual len: ", len(s.Mem.VirtualMemory), "; cap: ", cap(s.Mem.VirtualMemory))

		case "dump":
			fmt.Println("process dump:")
			for _, proc := range s.ReadyQ {
				fmt.Println(*proc)
			}

			for _, proc := range s.WaitingQ {
				fmt.Println(*proc)
			}

		case "exit":
			fmt.Println("exiting simulator")
			return

		default:
			fmt.Printf("ERROR: Command `%s` not found\n", args[0])
			break

		}
	}
}
