
package utils

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
	
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/sched"
)

// ShuffleInstructions : randomize the order of instructions
func ShuffleInstructions(vals [][]string) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	// We start at the end of the slice, inserting our random
	// values one at a time.
	for n := len(vals); n > 0; n-- {
		randIndex := r.Intn(n)
		// We swap the value at index n-1 and the random index
		// to move our randomly chosen value to the end of the
		// slice, and to move the value that was at n-1 into our
		// unshuffled portion of the slice.
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
	}
}

// LoadTemplate : load in template process and create process mutations off of it
func LoadTemplate(filename string, numOfProcesses int, processChan chan *sched.Process) error {

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file", err)
		return nil
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
		return err
	}

	// Randomize order of isntructions
	ShuffleInstructions(instructions)
	
	// program := code.Assemble(instructions)

	for i := 0; i < numOfProcesses; i++ {
		go sched.CreateRandomProcessFromTemplate(filename, instructions, processChan)
	}

	return nil
}