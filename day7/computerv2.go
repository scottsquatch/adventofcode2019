package main

// Computer represents a copmuter that runs IntCode programs
type ComputerV2 struct {
	input  <-chan int
	output chan<- int
	halt   chan<- bool
	memory []int
}

// NewComputer Create a new Copmuter
func NewComputerV2(in <-chan int, out chan<- int, halt chan<- bool) *ComputerV2 {
	return &ComputerV2{in, out, halt, make([]int, 0)}
}

// Run an IntCode program with the current computer
func (comp ComputerV2) Run(program *IntCodeProgram) {
	comp.memory = make([]int, len(program.programData))
	copy(comp.memory, program.programData)
	ip := 0
	halt := false
	for ip < len(comp.memory) {
		inst := newInstruction(comp.memory[ip])
		switch inst.op {
		case add:
			ip++
			a := getData(comp.memory[ip], inst.modes[0], comp.memory)
			ip++
			b := getData(comp.memory[ip], inst.modes[1], comp.memory)
			ip++
			c := getData(comp.memory[ip], inst.modes[2], comp.memory)

			comp.memory[c] = a + b

			ip++
		case mult:
			ip++
			a := getData(comp.memory[ip], inst.modes[0], comp.memory)
			ip++
			b := getData(comp.memory[ip], inst.modes[1], comp.memory)
			ip++
			c := getData(comp.memory[ip], inst.modes[2], comp.memory)

			comp.memory[c] = a * b

			ip++
		case store:
			x := <-comp.input
			ip++
			a := getData(comp.memory[ip], inst.modes[0], comp.memory)
			comp.memory[a] = x

			ip++
		case output:
			ip++
			a := getData(comp.memory[ip], inst.modes[0], comp.memory)
			comp.output <- a

			ip++
		case jumpft:
			ip++
			a := getData(comp.memory[ip], inst.modes[0], comp.memory)
			ip++
			b := getData(comp.memory[ip], inst.modes[1], comp.memory)
			if a != 0 {
				ip = b
			} else {
				ip++
			}
		case jumpff:
			ip++
			a := getData(comp.memory[ip], inst.modes[0], comp.memory)
			ip++
			b := getData(comp.memory[ip], inst.modes[1], comp.memory)
			if a == 0 {
				ip = b
			} else {
				ip++
			}
		case lt:
			ip++
			a := getData(comp.memory[ip], inst.modes[0], comp.memory)
			ip++
			b := getData(comp.memory[ip], inst.modes[1], comp.memory)
			ip++
			c := getData(comp.memory[ip], inst.modes[2], comp.memory)
			if a < b {
				comp.memory[c] = 1
			} else {
				comp.memory[c] = 0
			}
			ip++
		case eq:
			ip++
			a := getData(comp.memory[ip], inst.modes[0], comp.memory)
			ip++
			b := getData(comp.memory[ip], inst.modes[1], comp.memory)
			ip++
			c := getData(comp.memory[ip], inst.modes[2], comp.memory)
			if a == b {
				comp.memory[c] = 1
			} else {
				comp.memory[c] = 0
			}
			ip++
		case exit:
			comp.halt <- true
			halt = true
		default:
			panic("Unsupported opcode: " + string(inst.op))
		}

		if halt {
			break
		}
	}
	close(comp.output)
}
