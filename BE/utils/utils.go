
package utils

import (
	"bufio"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	
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

// ReadLine : read and parse line into 2d arrat
func ReadLine(reader *bufio.Reader) ([]string, error) {

	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.ReplaceAll(line, "\n", "")
	elements := strings.Split(line, " ")

	return elements, nil
}


func StrToIntArray(strArray []string) []int {

	ret := make([]int, 0)

	for _, str := range strArray {
		intElem, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println("error: invalid operand type")
		}

		ret = append(ret, intElem)
	}

	return ret
}