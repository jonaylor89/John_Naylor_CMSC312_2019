package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	// "time"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/sched"
)


func main() {

	// Message channel between main kernel and scheduler
	ch := make(chan *sched.Process, 1000)
	defer close(ch)

	s := sched.Scheduler{
		InMsg:     ch,
		ReadyQ: []*sched.Process{},
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

			numOfProc, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println("Could not get number of processes")
				break
			}

			if numOfProc <= 0 {
				fmt.Println("`load` number of processes must be postive")
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

			// Randomize order of isntructions
			sched.ShuffleInstructions(instructions)

			for i := 0; i < numOfProc; i++ {
				go sched.CreateRandomProcessFromTemplate(args[1], instructions, ch)
			}

		case "len":
			fmt.Println("ready: ", len(s.ReadyQ), "; waiting: ", len(s.WaitingQ))

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
			break

		}
	}
}
