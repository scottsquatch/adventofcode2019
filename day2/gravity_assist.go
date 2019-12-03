package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
)

func parseIntCodes(input string) []int {
	codes := strings.Split(strings.Trim(input, "\n"), ",")

	intCodes := make([]int, len(codes))

	for i, code := range codes {
		intCode, err := strconv.Atoi(code)
		if err == nil {
			intCodes[i] = intCode
		}
	}

	return intCodes
}

func runIntCodeProgram(computer *Computer, program *IntCodeProgram, noun int, verb int) []int {
	dataCopy := make([]int, len(program.programData))
	copy(dataCopy, program.programData)
	programCopy := NewIntCodeProgram(dataCopy)
	programCopy.programData[1] = noun
	programCopy.programData[2] = verb

	return computer.Run(programCopy).programData
}

func findInputs(computer *Computer, program *IntCodeProgram, target int) (int, int) {
	var lastOutput = 0

	// Based on analyzing the output, it looks like changing noun changes input a lot
	// but changing verb causes smaller changes
	var noun int = 0
	var verb int = 0

	for ; lastOutput < target; noun++ {
		lastOutput = runIntCodeProgram(computer, program, noun, verb)[0]
	}

	noun -= 2
	for ; lastOutput != target; verb++ {
		lastOutput = runIntCodeProgram(computer, program, noun, verb)[0]
	}

	return verb - 1, noun
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please enter the name of the file which contains the module masses when launching program")
		fmt.Println("Usage: day2 input.txt")
		return
	}

	computer := NewComputer()
	intCodes := parseIntCodes(fileutils.ReadFile(os.Args[1]))
	program := NewIntCodeProgram(intCodes)

	fmt.Printf("value of position 0 after program halts: %d\n", runIntCodeProgram(computer, program, 12, 2)[0])

	verb, noun := findInputs(computer, program, 19690720)
	fmt.Printf("noun %d, verb %d, 100 * noun + verb == %d\n", noun, verb, 100*noun+verb)
}
