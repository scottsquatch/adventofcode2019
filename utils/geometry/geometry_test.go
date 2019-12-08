package geometry

import (
	"reflect"
	"testing"
)

func TestNewLine(t *testing.T) {
	want := Line{Point{0, 0}, Point{1, 4}}
	got := NewLine(Point{0, 0}, Point{1, 4})

	if got.start != want.start {
		t.Errorf("got.start == %q, wanted %q", got.start, want.start)
	} else if got.end != want.end {
		t.Errorf("got.end == %q, wanted %q", got.end, want.end)
	}
}

func TestAddPoint(t *testing.T) {
	want := []Point{Point{0, 0}, Point{1, 4}}
	path := NewPath()
	path.AddPoint(Point{0, 0})
	path.AddPoint(Point{1, 4})
	got := path.points

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %q, wanted %q\n", got, want)
	}
}

func TestGetLines(t *testing.T) {
	in := &Path{[]Point{Point{0, 0}, Point{0, 4}, Point{2, 4}}}
	want := []Line{Line{Point{0, 0}, Point{0, 4}}, Line{Point{0, 4}, Point{2, 4}}}
	got := in.GetLines()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %q but wanted %q", got, want)
	}
}

func TestManhattanDistance(t *testing.T) {
	cases := []struct {
		in   [2]Point
		want int
	}{
		{[2]Point{Point{3, 3}, Point{0, 0}}, 6},
		{[2]Point{Point{0, 0}, Point{3, 3}}, 6},
	}

	for _, c := range cases {
		got := c.in[0].ManhattanDistance(c.in[1])
		if got != c.want {
			t.Errorf("%+v.ManhattanDistance(%+v) == %d, want %d", c.in[0], c.in[1], got, c.want)
		}
	}
}
