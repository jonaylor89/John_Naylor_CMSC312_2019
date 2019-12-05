package code

import (
	"bytes"
	"fmt"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/utils"
)

// Instructions : special type for list of instructions
type Instructions []byte

// Opcode : an opcode to an instruction
type Opcode byte

const (

	// CALCULATE : CPU operation
	CALC Opcode = iota

	// IO : i/o operation
	IO

	// FORK : fork a process
	FORK

	// ENTER : enter critical section
	ENTER 

	// EXIT : exit critical section
	EXIT

	// SEND : send packet to mailbox
	SEND

	// RECV : get packet from mailbox
	RECV

	// NOP : No operation
	NOP 
)

// Definition : definition of an instruction
type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	CALC: 	   {"CALC", []int{1}},
	IO:        {"IO", []int{1}},
	FORK:      {"FORK", []int{}},
	ENTER: 	   {"ENTER", []int{}},
	EXIT: 	   {"EXIT", []int{}},
	SEND: 	   {"SEND", []int{1}},
	RECV: 	   {"RECV", []int{1}},
	NOP:	   {"NOP", []int{}},
}

// Lookup : associate a opcode with its definition
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}
		

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

// Make : create instruction from opcode and operands
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1

	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1

	for i, o := range operands {
		width := def.OperandWidths[i]

		switch width {
		case 1:
			instruction[offset] = byte(o)
		}

		offset += width
	}

	return instruction
}

// ReadUint8 : read in an 8 bit unsigned integer
func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[0])
}

// ReadOperands : Get the operands of instructions
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

// Assemble : Assembly a 2 dimensions string array of opcode and operands into Instructions
func Assemble(instructions [][]string) Instructions { 

	program := Instructions{}
	var op []byte

	for _, ins := range instructions {
		switch ins[0] {
		case "CALC":
			op = Make(CALC, utils.StrToIntArray(ins[1:])...)
			break
		case "IO":
			op = Make(IO, utils.StrToIntArray(ins[1:])...)
			break
		case "ENTER":
			op = Make(ENTER, utils.StrToIntArray(ins[1:])...)
			break
		case "EXIT":
			op = Make(EXIT, utils.StrToIntArray(ins[1:])...)
			break
		case "FORK":
			op = Make(FORK, utils.StrToIntArray(ins[1:])...)
			break
		case "SEND":
			op = Make(SEND, utils.StrToIntArray(ins[1:])...)
			break
		case "RECV":
			op = Make(RECV, utils.StrToIntArray(ins[1:])...)
			break
		case "NOP":
			op = Make(NOP, utils.StrToIntArray(ins[1:])...)
			break
		default:
			op = Make(NOP, utils.StrToIntArray(ins[1:])...)
			break
		}

		program = append(program, op...)
	}

	return program
}
