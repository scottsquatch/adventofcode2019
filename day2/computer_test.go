package main

import "testing"

func TestComputerRun(t *testing.T) {
	comp := NewComputer()
	cases := []struct {
		in   *IntCodeProgram
		want *IntCodeProgram
	}{
		{NewIntCodeProgram([]int{1, 0, 0, 0, 99}), NewIntCodeProgram([]int{2, 0, 0, 0, 99})},
		{NewIntCodeProgram([]int{2, 3, 0, 3, 99}), NewIntCodeProgram([]int{2, 3, 0, 6, 99})},
		{NewIntCodeProgram([]int{2, 4, 4, 5, 99, 0}), NewIntCodeProgram([]int{2, 4, 4, 5, 99, 9801})},
		{NewIntCodeProgram([]int{1, 1, 1, 4, 99, 5, 6, 0, 99}), NewIntCodeProgram([]int{30, 1, 1, 4, 2, 5, 6, 0, 99})},
	}
	for _, c := range cases {
		got := comp.Run(c.in)
		for i, x := range got.programData {
			if x != c.want.programData[i] {
				t.Errorf("Computer.Run(%d) == %d, want %d", c.in, got, c.want.programData)
			}
		}
	}
}
