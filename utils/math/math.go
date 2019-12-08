package math

// Max return the max of two intergers
func Max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

// Min return the min of two intergers
func Min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

// Divmod calulates and return the division and mod operations on the two inputs
func Divmod(a int, b int) (div int, mod int) {
	return a / b, a % b
}
