package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
)

// IdealFuelCalculator calculates the fuel needed assuming that fuel is weightless
type IdealFuelCalculator struct{}

func fuelNeeded(mass uint64) uint64 {
	return mass/3 - 2
}

// Calculate calculates the fuel needed for the given masses
func (ifc IdealFuelCalculator) Calculate(masses []uint64) uint64 {
	var totalFuel uint64 = 0

	for _, mass := range masses {
		totalFuel += fuelNeeded(mass)
	}

	return totalFuel
}

// RealisticFuelCalculator calculates the fuel needed accounting for the weight of fuel
type RealisticFuelCalculator struct{}

func fuelNeededRecursive(mass uint64) uint64 {
	initialFuelNeeded := fuelNeeded(mass)

	// Be careful to check for overflow here
	if initialFuelNeeded <= 0 || initialFuelNeeded > mass {
		return 0
	}

	return initialFuelNeeded + fuelNeededRecursive(initialFuelNeeded)
}

// Calculate calculates the fuel needed for the given masses
func (rfc RealisticFuelCalculator) Calculate(masses []uint64) uint64 {
	var totalFuel uint64 = 0

	for _, mass := range masses {
		totalFuel += fuelNeededRecursive(mass)
	}

	return totalFuel
}

func readMasses(path string) []uint64 {
	fileContents := strings.Trim(fileutils.ReadFile(path), "\n")
	fileLines := strings.Split(fileContents, "\n")
	masses := make([]uint64, len(fileLines))

	for i, line := range fileLines {
		mass, err := strconv.Atoi(line)
		if err == nil {
			masses[i] = uint64(mass)
		}
	}

	return masses
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please enter the name of the file which contains the module masses when launching program")
		fmt.Println("Usage: day1 input.txt")
		return
	}

	masses := readMasses(os.Args[1])
	idealCalc := IdealFuelCalculator{}
	realCalc := RealisticFuelCalculator{}

	fmt.Printf("Total required fuel (assuming 0 fuel mass) is %d\n", idealCalc.Calculate(masses))
	fmt.Printf("Total required fuel (accounting for fuel mass) is %d\n", realCalc.Calculate(masses))
}
