package fileutils

// Utility methods pertaining to file input/output

import "io/ioutil"

// Helper method to throw error if it exists
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

// Simple method to read a file into memory and output its contents
func ReadFile(path string) string {
	dat, err := ioutil.ReadFile(path)
	checkError(err)

	return string(dat)
}