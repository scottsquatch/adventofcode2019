package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/scottsquatch/adventofcode2019/utils/math"
)

// IntCodeProgram represents a program to be run by Computer
type IntCodeProgram struct {
	programData []int
}

// Computer represents a copmuter that runs IntCode programs
type Computer struct {
}

// NewIntCodeProgram Create a new IntCode Program
func NewIntCodeProgram(intCodes []int) *IntCodeProgram {
	return &IntCodeProgram{intCodes}
}

// NewComputer Create a new Copmuter
func NewComputer() *Computer {
	return &Computer{}
}

// Run an IntCode program with the current computer
func (comp Computer) Run(program *IntCodeProgram) *IntCodeProgram {
	currentProgram := IntCodeProgram{program.programData}
	scanner := bufio.NewScanner(os.Stdin)
	ip := 0
	halt := false
	for ip < len(currentProgram.programData) {
		inst := newInstruction(currentProgram.programData[ip])

		switch inst.op {
		case add:
			ip++
			a := getData(currentProgram.programData[ip], inst.modes[0], currentProgram.programData)
			ip++
			b := getData(currentProgram.programData[ip], inst.modes[1], currentProgram.programData)
			ip++
			c := getData(currentProgram.programData[ip], inst.modes[2], currentProgram.programData)

			currentProgram.programData[c] = a + b

			ip++
		case mult:
			ip++
			a := getData(currentProgram.programData[ip], inst.modes[0], currentProgram.programData)
			ip++
			b := getData(currentProgram.programData[ip], inst.modes[1], currentProgram.programData)
			ip++
			c := getData(currentProgram.programData[ip], inst.modes[2], currentProgram.programData)

			currentProgram.programData[c] = a * b

			ip++
		case store:
			fmt.Print("input: ")
			scanner.Scan()
			x, _ := strconv.Atoi(scanner.Text())
			ip++
			a := getData(currentProgram.programData[ip], inst.modes[0], currentProgram.programData)
			currentProgram.programData[a] = x

			ip++
		case output:
			ip++
			a := getData(currentProgram.programData[ip], inst.modes[0], currentProgram.programData)
			fmt.Println(a)

			ip++
		case jumpft:
			ip++
			a := getData(currentProgram.programData[ip], inst.modes[0], currentProgram.programData)
			ip++
			b := getData(currentProgram.programData[ip], inst.modes[1], currentProgram.programData)
			if a != 0 {
				ip = b
			} else {
				ip++
			}
		case jumpff:
			ip++
			a := getData(currentProgram.programData[ip], inst.modes[0], currentProgram.programData)
			ip++
			b := getData(currentProgram.programData[ip], inst.modes[1], currentProgram.programData)
			if a == 0 {
				ip = b
			} else {
				ip++
			}
		case lt:
			ip++
			a := getData(currentProgram.programData[ip], inst.modes[0], currentProgram.programData)
			ip++
			b := getData(currentProgram.programData[ip], inst.modes[1], currentProgram.programData)
			ip++
			c := getData(currentProgram.programData[ip], inst.modes[2], currentProgram.programData)
			if a < b {
				currentProgram.programData[c] = 1
			} else {
				currentProgram.programData[c] = 0
			}
			ip++
		case eq:
			ip++
			a := getData(currentProgram.programData[ip], inst.modes[0], currentProgram.programData)
			ip++
			b := getData(currentProgram.programData[ip], inst.modes[1], currentProgram.programData)
			ip++
			c := getData(currentProgram.programData[ip], inst.modes[2], currentProgram.programData)
			if a == b {
				currentProgram.programData[c] = 1
			} else {
				currentProgram.programData[c] = 0
			}
			ip++
		case exit:
			halt = true
		default:
			panic("Unsupported opcode: " + string(inst.op))
		}

		if halt {
			break
		}
	}

	return &currentProgram
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
	exit   opCode = 99
)

type parameterMode int

const (
	position  parameterMode = 0
	immediate parameterMode = 1
)

type instruction struct {
	op    opCode
	modes []parameterMode
}

func getOp(op int) opCode {
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
	case 99:
		return exit
	default:
		panic("Op code of " + strconv.Itoa(op) + " is not valid")
	}
}

func getparameterMode(i int) parameterMode {
	switch i {
	case 0:
		return position
	case 1:
		return immediate
	default:
		panic("Instruction mode of " + strconv.Itoa(i) + " is not valid")
	}
}

func newInstruction(i int) instruction {
	str := strconv.Itoa(i)
	rawOp, _ := strconv.Atoi(str[math.Max(len(str)-2, 0):])
	op := getOp(rawOp)

	var modes []parameterMode
	switch op {
	case add, mult, lt, eq:
		modes = make([]parameterMode, 3)
	case jumpff, jumpft:
		modes = make([]parameterMode, 2)
	case store, output:
		modes = make([]parameterMode, 1)
	case exit:
		modes = make([]parameterMode, 0)
	}

	for i := 0; i < len(modes); i++ {
		m := position
		idx := len(str) - 3 - i
		if idx >= 0 {
			rawMode, _ := strconv.Atoi(string(str[idx]))
			m = getparameterMode(rawMode)
		}
		modes[i] = m
	}

	switch op {
	case add, mult, lt, eq:
		modes[2] = immediate
	case store:
		modes[0] = immediate
	}

	return instruction{op, modes}
}

func getData(val int, m parameterMode, programData []int) int {
	switch m {
	case immediate:
		return val
	case position:
		return programData[val]
	default:
		panic("Parameter mode of " + string(m) + " is not valid")
	}
}
