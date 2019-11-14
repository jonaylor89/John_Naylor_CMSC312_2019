package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	// "time"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/kernel"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/memory"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/cpu"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/server"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/config"
)

func main() {

	conf := config.ReadConfig("config.yml")

	// Message channel to kernel
	ch := make(chan *kernel.Process, conf.ProcChanSize)
	defer close(ch)

	cpu := &cpu.CPU{ 
		TotalCycles: 0, 
		Speed: conf.CPU.ClockSpeed,
	}

	mem := &memory.Memory{
		PageSize: conf.Memory.PageSize,
		TotalRam: conf.Memory.TotalRam,
		PageTable: make(map[int]int),
		VirtualMemory: make([]*memory.Page, 0),
		PhysicalMemory: make([]*memory.Page, 0, conf.Memory.TotalRam / conf.Memory.PageSize),
	}

	k := &kernel.Kernel{
		CPU: 	  cpu,
		Mem: 	  mem,
		InMsg:    ch,
		ReadyQ:   []*kernel.Process{},
		WaitingQ: []*kernel.Process{},
		MinimumFreeFrames: conf.MinimumFreeFrames,
		TimeQuantum: conf.Sched.TimeQuantum,
		Mailboxes: []chan byte {
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

	// Run the scheduler
	go k.RunRoundRobin()
	// go k.RunFirstComeFirst
	
	go server.StartServer(k)

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

			err = kernel.LoadTemplate(filename, numOfProc, ch)
			if err != nil {
				fmt.Println("`load` error loading process template", err)
				break
			}

		case "proc":
			fmt.Println("ready: ", len(k.ReadyQ), "; waiting: ", len(k.WaitingQ), "; sending: ", len(ch))

		case "mem":
			fmt.Println("Physical len: ", len(k.Mem.PhysicalMemory), "; cap: ", cap(k.Mem.PhysicalMemory))
			fmt.Println("Virtual len: ", len(k.Mem.VirtualMemory), "; cap: ", cap(k.Mem.VirtualMemory))

		case "dump":
			fmt.Println("process dump:")
			for _, proc := range k.ReadyQ {
				fmt.Println(*proc)
			}

			for _, proc := range k.WaitingQ {
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
