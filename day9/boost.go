package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
)

func runPart1() {
	program := NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))

	input := make(chan int64)
	output := make(chan int64)
	halt := make(chan bool)
	vm := NewComputer(input, output)
	go func() {
		input <- 1
	}()

	wg := sync.WaitGroup{}

	go func() {
		defer func() {
			wg.Done()
		}()
		wg.Add(1)
		for {
			select {
			case x := <-output:
				fmt.Printf("Boost code: %d\n", x)
			case <-halt:
				return
			}
		}
	}()
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()

		vm.Run(program)
		halt <- true
	}()

	wg.Wait()
}

func runPart2() {
	program := NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))

	input := make(chan int64)
	output := make(chan int64)
	halt := make(chan bool)
	vm := NewComputer(input, output)
	go func() {
		input <- 2
	}()

	wg := sync.WaitGroup{}

	go func() {
		defer func() {
			wg.Done()
		}()
		wg.Add(1)
		for {
			select {
			case x := <-output:
				fmt.Printf("Distress Signal: %d\n", x)
			case <-halt:
				return
			}
		}
	}()
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()

		vm.Run(program)
		halt <- true
	}()

	wg.Wait()
}

func main() {
	runPart1()
	runPart2()
}
