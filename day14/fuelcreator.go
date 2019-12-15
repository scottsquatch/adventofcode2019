package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
)

type chemicalRequirement struct {
	chemical string
	quantity int
}

type chemicalRequirements struct {
	components []chemicalRequirement
	quantity   int
}

func parseChemicalRequirement(req string) chemicalRequirement {
	comp := strings.Split(req, " ")
	quant, _ := strconv.Atoi(comp[0])

	return chemicalRequirement{comp[1], quant}
}

func parseChemicalRequirements(reqs string) map[string]chemicalRequirements {
	chemicals := make(map[string]chemicalRequirements)

	for _, str := range strings.Split(fileutils.ReadFile(os.Args[1]), "\n") {
		combination := strings.Split(str, " => ")
		components := []chemicalRequirement{}
		for _, compstr := range strings.Split(combination[0], ", ") {
			components = append(components, parseChemicalRequirement(compstr))
		}
		result := parseChemicalRequirement(combination[1])
		chemicals[result.chemical] = chemicalRequirements{components, result.quantity}
	}

	return chemicals
}

func getOresNeeded(chemicals map[string]chemicalRequirements, fuelAmount int) int64 {
	var ores int64 = 0
	stack := []chemicalRequirement{chemicalRequirement{"FUEL", fuelAmount}}
	var chem chemicalRequirement
	inventory := make(map[string]int)
	for len(stack) > 0 {
		chem, stack = stack[len(stack)-1], stack[:len(stack)-1]

		reaction := chemicals[chem.chemical]
		need := chem.quantity
		have := inventory[chem.chemical]

		if have >= need {
			inventory[chem.chemical] -= need
		} else {
			need -= have
			get := reaction.quantity
			num := int(math.Ceil(float64(need) / float64(get)))
			for _, req := range chemicals[chem.chemical].components {
				make := num * req.quantity
				if req.chemical == "ORE" {
					ores += int64(make)
				} else {
					stack = append(stack, chemicalRequirement{req.chemical, make})
				}
			}
			inventory[chem.chemical] += num*reaction.quantity - chem.quantity
		}
	}

	return ores
}

func runPartA() map[string]chemicalRequirements {
	chemicals := parseChemicalRequirements(os.Args[1])

	ores := getOresNeeded(chemicals, 1)

	fmt.Printf("Total number of ores needed for 1 fuel is %d\n", ores)

	return chemicals
}

func runPartB(chemicals map[string]chemicalRequirements) {
	var target, last int64 = 1000000000000, -1
	fuel := 1
	var low, high int
	for last < target {
		low = fuel
		last = getOresNeeded(chemicals, fuel)
		fuel = fuel * 2
	}

	high = low
	low /= 2

	lastFuel := fuel - 1
	for last != target && lastFuel != fuel {
		if last > target {
			high = fuel
		} else {
			low = fuel
		}

		lastFuel = fuel
		fuel = (high + low) / 2

		last = getOresNeeded(chemicals, fuel)
	}

	fmt.Printf("Amount of fuel that can be produced from %d ores is %d\n", target, fuel)
}

func main() {
	chemicals := runPartA()
	runPartB(chemicals)
}
