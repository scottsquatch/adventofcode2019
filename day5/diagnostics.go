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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please enter the name of the file which contains the module masses when launching program")
		fmt.Println("Usage: day5 input.txt")
		return
	}

	computer := NewComputer()
	intCodes := parseIntCodes(fileutils.ReadFile(os.Args[1]))
	program := NewIntCodeProgram(intCodes)

	fmt.Println("Starting diagnostic")

	computer.Run(program)

}
