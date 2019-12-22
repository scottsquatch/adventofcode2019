package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/scottsquatch/adventofcode2019/common"
	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
)

func runPartA() {
	in, out, done := make(chan int64), make(chan int64), make(chan bool)
	software := common.NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))
	asciirobot := common.NewComputer(in, out)

	go func() {
		asciirobot.Run(software)
		close(out)
	}()

	go func() {
		for c := range out {
			var str string
			if c >= 0 && c <= 127 { // ASCII char
				str = string(c)
			} else {
				str = strconv.FormatInt(c, 10)
			}

			fmt.Print(str)
		}
		done <- true
		fmt.Println()
	}()

	go func() {
		// Determined cases by realizing that robot jumps 4 tiles, thus
		// 	Cases:
		// 		- 1 Square hole -> Start jump 3 squares early
		// 		- 2 square hole -> start jump 2 squares early
		// 		- 3 square hole -> start jump 1 square early
		// 		- 4 square hole -> Die
		//      - next square is hole -> jump and hope we make it
		// Which can be written as (A*B*!C*D)+(A*!B*!C*D)+(!A*!B*!C*D)+!A, which reduces to (!C*D)(A*B+!B)+!A
		commands := []string{"OR A T\n", "AND B T\n", "NOT B J\n", "OR T J\n", "NOT C T\n", "AND D T\n", "AND T J\n", "NOT A T\n", "OR T J\n", "WALK\n"}
		for _, cmd := range commands {
			for _, c := range cmd {
				in <- int64(c)
			}
		}
	}()

	<-done
}

func runPartB() {
	//D && (!A || (H && (!C || (!B && (!E || !F || !G)))
	run([]string{"NOT G T\n", "NOT E J\n", "OR T J\n", "NOT F T\n", "OR J T\n", "NOT B J\n", "AND J T\n", "NOT C J\n", "OR T J\n", "AND H J\n", "NOT A T\n", "OR T J\n", "AND D J\n", "RUN\n"})
}

func run(commands []string) bool {
	in, out, done := make(chan int64), make(chan int64), make(chan bool)
	software := common.NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))
	asciirobot := common.NewComputer(in, out)

	go func() {
		asciirobot.Run(software)
		close(out)
	}()

	go func() {
		foundSolution := false
		for c := range out {
			var str string
			if c >= 0 && c <= 127 { // ASCII char
				str = string(c)
			} else {
				str = strconv.FormatInt(c, 10)
				foundSolution = true
			}

			fmt.Print(str)
		}
		fmt.Println()
		done <- foundSolution
	}()

	go func() {
		for _, cmd := range commands {
			for _, c := range cmd {
				in <- int64(c)
			}
		}
	}()

	return <-done
}

func main() {
	runPartA()
	runPartB()
}
