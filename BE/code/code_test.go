package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{CALC, []int{130}, []byte{byte(CALC), 130}},
		{IO, []int{255}, []byte{byte(IO), 255}},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want=%d, got=%d", len(tt.expected), len(instruction))
		}

		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d", i, b, instruction[i])
			}
		}
	}
}

func TestInstructionString(t *testing.T) {

	instructions := []Instructions{
		Make(CALC, 1),
		Make(IO, 2),
	}

	expected := `0000 CALC 1
0002 IO 2
`

	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted.\nwant=%q\ngot=%q",
			expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{CALC, []int{130}, 1},
		{IO, []int{255}, 1},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("n wrong. want=%d, got=%d", tt.bytesRead, n)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("operand wrong. want=%d, got=%d", want, operandsRead[i])
			}
		}
	}
}

func TestAssemble(t *testing.T) {
	instructions := [][]string{
		[]string{"CALC", "32"},
		[]string{"IO", "14"},
	}

	expected := `0000 CALC 32
0002 IO 14
`

	program := Assemble(instructions)

	if program.String() != expected {
		t.Errorf("instructions wrongly formatted.\nwant=%q\ngot=%q",
			expected, program.String())
	}
}