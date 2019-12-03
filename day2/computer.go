package main

// IntCodeProgram represents a program to be run by Computer
type IntCodeProgram struct {
	programData []int
}

// Computer represents a copmuter that runs IntCode programs
type Computer struct{}

// NewIntCodeProgram Create a new IntCode Program
func NewIntCodeProgram(intCodes []int) *IntCodeProgram {
	return &IntCodeProgram{intCodes}
}

// NewComputer Create a new Copmuter
func NewComputer() *Computer {
	return &Computer{}
}

const add = 1
const mult = 2
const exit = 99

// Run an IntCode program with the current computer
func (comp Computer) Run(program *IntCodeProgram) *IntCodeProgram {
	currentProgram := IntCodeProgram{program.programData}

	for pc := 0; pc < len(currentProgram.programData); pc += 4 {
		opcode := currentProgram.programData[pc]
		if opcode == exit {
			break
		}

		posA := currentProgram.programData[pc+1]
		posB := currentProgram.programData[pc+2]
		posC := currentProgram.programData[pc+3]

		var c int
		if opcode == add {
			c = currentProgram.programData[posA] + currentProgram.programData[posB]
		} else {
			c = currentProgram.programData[posA] * currentProgram.programData[posB]
		}

		currentProgram.programData[posC] = c
	}

	return &currentProgram
}
