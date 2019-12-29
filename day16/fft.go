package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
)

func abs(a int) int {
	if a >= 0 {
		return a
	}

	return -1 * a
}

func runPartA(input []int) {
	// Looking at the example, all values before an index are 0
	phases := 100
	output := make([]int, len(input))
	in := make([]int, len(input))
	copy(in, input)
	for p := 0; p < phases; p++ {
		for i := 0; i < len(output); i++ {
			c := 1
			s := 0
			j := i
			for j < len(input) {
				for k := 0; k < i+1 && j < len(input); k++ {
					s += c * in[j]
					j++
				}
				c *= -1
				j += i + 1
			}
			output[i] = abs(s) % 10
		}
		copy(in, output)
	}

	fmt.Printf("First 8 digits of output after %d phases is %v\n", phases, in[:8])
}

func runPartB(input []int) {
	// from the examples, digits from len/2 to end are can be calculated from previous results
	offset := make([]string, 7)
	for i := range offset {
		offset[i] = strconv.FormatInt(int64(input[i]), 10)
	}

	offsetval, _ := strconv.Atoi(strings.Join(offset, ""))
	repeated := 10000
	phases := 100
	N := len(input)
	partial := make([]int, N*repeated-offsetval)
	for i, j := 0, offsetval; i < len(partial); i, j = i+1, j+1 {
		partial[i] = input[j%N]
	}
	out := make([]int, len(partial))
	for p := 0; p < phases; p++ {
		// From observation out[i] = sum(in[i:len(in)-1]), given i > len(in)/2
		sum := 0
		for i := len(partial) - 1; i >= 0; i-- {
			sum += partial[i]
			out[i] = sum % 10
		}
		copy(partial, out)
	}

	fmt.Printf("Message after %d phases is %v\n", phases, partial[:8])
}

func main() {
	signalstr := strings.Trim(fileutils.ReadFile(os.Args[1]), " \r\n")
	input := make([]int, len(signalstr))
	for i, r := range signalstr {
		input[i] = int(r - '0')
	}
	runPartA(input)
	runPartB(input)
}
