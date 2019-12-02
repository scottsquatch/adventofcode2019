package fileutils

import "testing"

func TestReadFile(t *testing.T) {
	expected := "Hello, World!"
	got := ReadFile("testdata.txt")

	if got != expected {
		t.Errorf("ReadFile(\"testdata.txt\") == %q, want %q", got, expected)
	}
}
