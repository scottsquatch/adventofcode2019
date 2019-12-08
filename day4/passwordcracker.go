package main

import (
	"fmt"
)

type passwordValidator func([]int) bool

func isValidPassword(digits []int) bool {
	if len(digits) != 6 {
		return false
	}

	twoSame := false
	prev := digits[0]
	for _, d := range digits[1:] {
		if prev == d {
			twoSame = true
		}
		if d < prev {
			return false
		}

		prev = d
	}

	return twoSame
}

func numValidPasswords(start int, end int, valid passwordValidator) int {
	numValid := 0

	for pass := start; pass <= end; pass++ {
		if valid(toDigits(pass)) {
			numValid++
		}
	}

	return numValid
}

func isValidPassword2(digits []int) bool {
	if len(digits) != 6 {
		return false
	}

	frequencies := map[int]int{}
	prev := digits[0]
	frequencies[prev] = 1
	for _, d := range digits[1:] {
		if d < prev {
			return false
		}

		frequencies[d]++

		prev = d
	}

	for _, val := range frequencies {
		if val == 2 {
			return true
		}
	}

	return false
}

func toDigits(n int) []int {
	// Find largest base power of 10
	divisor := 1
	for divisor < n {
		divisor *= 10
	}

	divisor /= 10

	var digits []int
	var digit int
	i := n
	for i > 0 {
		digit, i = divmod(i, divisor)
		digits = append(digits, digit)
		divisor /= 10
	}

	return digits
}

func divmod(a int, b int) (div int, mod int) {
	return a / b, a % b
}

func main() {
	fmt.Printf("Number of valid passwords in given range %d\n", numValidPasswords(123257, 647015, isValidPassword))
	fmt.Printf("Number of valid passwords in given range (part 2) %d\n", numValidPasswords(123257, 647015, isValidPassword2))
}
