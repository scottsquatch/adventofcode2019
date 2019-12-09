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

func getPermutations(nums []int) [][]int {
	var perms [][]int
	var rc func([]int, int)
	rc = func(a []int, k int) {
		if k == len(a) {
			perms = append(perms, append([]int{}, a...))
		} else {
			for i := k; i < len(nums); i++ {
				a[k], a[i] = a[i], a[k]
				rc(a, k+1)
				a[k], a[i] = a[i], a[k]
			}
		}
	}

	rc(nums, 0)

	return perms
}

func getThrusterSignal(program *IntCodeProgram, phaseSequence []int) int {
	in := make(chan int, 1)
	out := make(chan int)
	halt := make(chan bool)
	computer := NewComputerV2(in, out, halt)

	outPhase := 0
	for _, phase := range phaseSequence {
		go computer.Run(program)
		in <- phase
		in <- outPhase
		outPhase = <-out
	}

	return outPhase
}

func findMaxThrusterSignal(program *IntCodeProgram) int {
	max := 0 // Assuming that signal is greater than 0

	phaseSequences := getPermutations([]int{4, 3, 2, 1, 0})

	for _, phaseSequence := range phaseSequences {
		thrusterSignal := getThrusterSignal(program, phaseSequence)
		if thrusterSignal > max {
			max = thrusterSignal
		}
	}

	return max
}

func getThrusterSignalWithFeedback(program *IntCodeProgram, phaseSequence []int) int {
	outChannels := make([]chan int, len(phaseSequence))
	inChannels := make([]chan int, len(phaseSequence))
	haltChannels := make([]chan bool, len(phaseSequence))
	computers := make([]*ComputerV2, len(phaseSequence))
	for i := 0; i < len(phaseSequence); i++ {
		outChannels[i] = make(chan int)
		inChannels[i] = make(chan int, 1)
		haltChannels[i] = make(chan bool)
		computers[i] = NewComputerV2(inChannels[i], outChannels[i], haltChannels[i])
	}

	for i, p := range phaseSequence {
		inChannels[i] <- p
		go computers[i].Run(program)
	}
	inChannels[0] <- 0

	out := -1
	done := make(chan bool)
	for i, c := range computers {
		go func(i int, c *ComputerV2) {
			for {
				select {
				case x := <-outChannels[i]:
					if i == 4 {
						out = x
						inChannels[0] <- x
					} else {
						inChannels[i+1] <- x
					}
				case _ = <-haltChannels[i]:
					if i == 4 {
						done <- true
					}
					break
				}
			}
		}(i, c)
	}

	<-done
	return out
}

func findMaxThrusterSignalWithFeedback(program *IntCodeProgram) int {
	max := 0 // Assuming that signal is greater than 0

	phaseSequences := getPermutations([]int{9, 8, 7, 6, 5})

	for _, phaseSequence := range phaseSequences {
		thrusterSignal := getThrusterSignalWithFeedback(program, phaseSequence)
		if thrusterSignal > max {
			max = thrusterSignal
		}
	}

	return max
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please enter the name of the file which contains the module masses when launching program")
		fmt.Println("Usage: day7 input.txt")
		return
	}

	intCodes := parseIntCodes(fileutils.ReadFile(os.Args[1]))
	program := NewIntCodeProgram(intCodes)

	fmt.Printf("Max Thruster Signal %d\n", findMaxThrusterSignal(program))
	fmt.Printf("Max Thruster Signal %d\n", findMaxThrusterSignalWithFeedback(program))
}
