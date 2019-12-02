package main

import "testing"

func TestIdealFuelCalculation(t *testing.T) {
	calc := IdealFuelCalculator{}
	cases := []struct {
		in   []uint64
		want uint64
	}{
		{[]uint64{12}, 2},
		{[]uint64{14}, 2},
		{[]uint64{1969}, 654},
		{[]uint64{100756}, 33583},
		{[]uint64{14, 1969}, 656},
	}
	for _, c := range cases {
		got := calc.Calculate(c.in)
		if got != c.want {
			t.Errorf("IdealFuelCalculator.Calculate(%d) == %d, want %d", c.in, got, c.want)
		}
	}
}

func TestRealisticFuelCalculation(t *testing.T) {
	calc := RealisticFuelCalculator{}
	cases := []struct {
		in   []uint64
		want uint64
	}{
		{[]uint64{14}, 2},
		{[]uint64{1969}, 966},
		{[]uint64{100756}, 50346},
		{[]uint64{14, 1969}, 968},
	}
	for _, c := range cases {
		got := calc.Calculate(c.in)
		if got != c.want {
			t.Errorf("RealisticFuelCalculator.Calculate(%d) == %d, want %d", c.in, got, c.want)
		}
	}
}
