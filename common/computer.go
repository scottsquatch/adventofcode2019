package common

import (
	"strconv"
	"strings"
)

// IntCodeProgram represents a program to be run by Computer
type IntCodeProgram struct {
	programData []int
}

// Computer represents a copmuter that runs IntCode programs
type Computer struct {
	input        <-chan int64
	output       chan<- int64
	memory       []int64
	relativeBase int64
	ip           int64
}

// NewIntCodeProgram Create a new IntCode Program
func NewIntCodeProgram(str string) *IntCodeProgram {
	codes := strings.Split(str, ",")
	intCodes := make([]int, len(codes))
	for i, c := range codes {
		intCodes[i], _ = strconv.Atoi(c)
	}
	return &IntCodeProgram{intCodes}
}

// NewComputer Create a new Copmuter
func NewComputer(in <-chan int64, out chan<- int64) *Computer {
	return &Computer{in, out, make([]int64, 2048), 0, 0}
}

// Run an IntCode program with the current computer
func (comp *Computer) Run(program *IntCodeProgram) {
	// copy
	for i, n := range program.programData {
		comp.memory[i] = int64(n)
	}

	comp.ip = 0
	halt := false
	for !halt {
		inst := newInstruction(comp.memory[comp.ip])
		switch inst.op {
		case add:
			a := comp.getData(comp.memory[comp.ip+1], inst.modes[0])
			b := comp.getData(comp.memory[comp.ip+2], inst.modes[1])
			c := comp.getDataWritingParameter(comp.memory[comp.ip+3], inst.modes[2])

			comp.memory[c] = a + b

			comp.ip += 4
		case mult:
			a := comp.getData(comp.memory[comp.ip+1], inst.modes[0])
			b := comp.getData(comp.memory[comp.ip+2], inst.modes[1])
			c := comp.getDataWritingParameter(comp.memory[comp.ip+3], inst.modes[2])

			comp.memory[c] = a * b

			comp.ip += 4
		case store:
			x := <-comp.input
			a := comp.getDataWritingParameter(comp.memory[comp.ip+1], inst.modes[0])
			comp.memory[a] = x

			comp.ip += 2
		case output:
			a := comp.getData(comp.memory[comp.ip+1], inst.modes[0])
			comp.output <- a

			comp.ip += 2
		case jumpft:
			a := comp.getData(comp.memory[comp.ip+1], inst.modes[0])
			b := comp.getData(comp.memory[comp.ip+2], inst.modes[1])
			if a != 0 {
				comp.ip = b
			} else {
				comp.ip += 3
			}
		case jumpff:
			a := comp.getData(comp.memory[comp.ip+1], inst.modes[0])
			b := comp.getData(comp.memory[comp.ip+2], inst.modes[1])
			if a == 0 {
				comp.ip = b
			} else {
				comp.ip += 3
			}
		case lt:
			a := comp.getData(comp.memory[comp.ip+1], inst.modes[0])
			b := comp.getData(comp.memory[comp.ip+2], inst.modes[1])
			c := comp.getDataWritingParameter(comp.memory[comp.ip+3], inst.modes[2])
			if a < b {
				comp.memory[c] = 1
			} else {
				comp.memory[c] = 0
			}
			comp.ip += 4
		case eq:
			a := comp.getData(comp.memory[comp.ip+1], inst.modes[0])
			b := comp.getData(comp.memory[comp.ip+2], inst.modes[1])
			c := comp.getDataWritingParameter(comp.memory[comp.ip+3], inst.modes[2])
			if a == b {
				comp.memory[c] = 1
			} else {
				comp.memory[c] = 0
			}
			comp.ip += 4
		case adj:
			a := comp.getData(comp.memory[comp.ip+1], inst.modes[0])

			comp.relativeBase += a

			comp.ip += 2
		case exit:
			halt = true
		default:
			panic("Unsupported opcode: " + string(inst.op))
		}
	}
}

// Op Codes
type opCode int

const (
	add    opCode = 1
	mult   opCode = 2
	store  opCode = 3
	output opCode = 4
	jumpft opCode = 5
	jumpff opCode = 6
	lt     opCode = 7
	eq     opCode = 8
	adj    opCode = 9
	exit   opCode = 99
)

type parameterMode int

const (
	position  parameterMode = 0
	immediate parameterMode = 1
	relative  parameterMode = 2
)

type instruction struct {
	op    opCode
	modes []parameterMode
}

func getOp(op int64) opCode {
	switch op {
	case 1:
		return add
	case 2:
		return mult
	case 3:
		return store
	case 4:
		return output
	case 5:
		return jumpft
	case 6:
		return jumpff
	case 7:
		return lt
	case 8:
		return eq
	case 9:
		return adj
	case 99:
		return exit
	default:
		panic("Op code of " + strconv.FormatInt(op, 10) + " is not valid")
	}
}

func getparameterMode(i int64) parameterMode {
	switch i {
	case 0:
		return position
	case 1:
		return immediate
	case 2:
		return relative
	default:
		panic("Instruction mode of " + strconv.FormatInt(i, 10) + " is not valid")
	}
}

func newInstruction(i int64) instruction {
	op := getOp(i % 100)

	var modes []parameterMode
	switch op {
	case add, mult, lt, eq:
		modes = make([]parameterMode, 3)
	case jumpff, jumpft:
		modes = make([]parameterMode, 2)
	case store, output, adj:
		modes = make([]parameterMode, 1)
	case exit:
		modes = make([]parameterMode, 0)
	}

	n := i / 100
	for j := 0; j < len(modes); j++ {
		m := parameterMode(n % 10)
		modes[j] = m
		n /= 10
	}

	return instruction{op, modes}
}

func (comp *Computer) getData(val int64, m parameterMode) int64 {
	switch m {
	case position:
		return comp.memory[val]
	case immediate:
		return val
	case relative:
		return comp.memory[comp.relativeBase+val]
	default:
		panic("Parameter mode of " + string(m) + " is not valid")
	}
}

func (comp *Computer) getDataWritingParameter(val int64, m parameterMode) int64 {
	switch m {
	case position, immediate:
		return val
	case relative:
		return comp.relativeBase + val
	default:
		panic("Parameter mode of " + string(m) + " is not valid")
	}
}
