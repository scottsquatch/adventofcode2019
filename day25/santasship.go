package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/scottsquatch/adventofcode2019/common"
	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
)

func runPartA() {
	in, out := make(chan int64), make(chan int64)
	program := common.NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))
	droid := common.NewComputer(in, out)
	reader := bufio.NewReader(os.Stdin)
	go func() {
		droid.Run(program)
		close(out)
	}()

	go func() {
		for _, dir := range strings.Split(directions, "\r\n") {
			for _, c := range dir {
				in <- int64(c)
			}
			in <- int64('\n')
		}
		// Brute force all combinations
		for _, el := range items {
			cmd := "drop " + el + "\n"
			for _, c := range cmd {
				in <- int64(c)
			}
		}
		for _, comb := range combinations(items) {
			for _, cmb := range comb {
				cmd := "take " + cmb + "\n"
				for _, c := range cmd {
					in <- int64(c)
				}
			}

			for _, c := range "east\n" {
				in <- int64(c)
			}

			// go back to initial state
			for _, cmb := range comb {
				cmd := "drop " + cmb + "\n"
				for _, c := range cmd {
					in <- int64(c)
				}
			}
		}
		for {
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\r\n", "", -1)
			for _, c := range text {
				in <- int64(c)
			}
			in <- int64('\n')
		}
	}()

	for o := range out {
		fmt.Print(string(o))
	}
}

func combinations(a []string) [][]string {
	max := 1
	for i := 0; i < len(a); i++ {
		max *= 2
	}

	combinations := [][]string{}
	for i := 1; i < max; i++ {
		comb := []string{}
		for j, num := range a {
			if i&(1<<j) == (1 << j) {
				comb = append(comb, num)
			}
		}
		combinations = append(combinations, comb)
	}

	return combinations
}

func runPartB() {

}

func main() {
	runPartA()
	runPartB()
}

var items = []string{"mug", "easter egg", "asterisk", "jam", "klein bottle", "tambourine", "cake", "polygon"}
var directions = `north
west
take mug
west
take easter egg
east
east
south
south
take asterisk
south
west
north
take jam
south
east
north
east
take klein bottle
south
west
take tambourine
west
take cake
east
south
east
take polygon
north`
