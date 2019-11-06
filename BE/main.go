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

	cpu := sched.CPU{ 
		TotalCycles: 0, 
		Speed: 10,
	}

	ram := memory.RAM{
		// frames: make([]*memory.Page, 0, memory.FrameLength),
	}

	s := sched.Scheduler{
		CPU: 	  cpu,
		RAM: 	  ram,
		InMsg:    ch,
		ReadyQ:   []*sched.Process{},
		WaitingQ: []*sched.Process{},
	}

	// Run the scheduler
	go s.RunRoundRobin()

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

		case "len":
			fmt.Println("ready: ", len(s.ReadyQ), "; waiting: ", len(s.WaitingQ), "; sending: ", len(ch))

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
