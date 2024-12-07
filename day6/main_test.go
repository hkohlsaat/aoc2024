package main

import (
	"fmt"
	"testing"
)

const input = `....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...`

var planMap = [][]rune{
	{'.', '.', '.', '.', '#', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '#'},
	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '#', '.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '.', '#', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
	{'.', '#', '.', '.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '.', '.', '#', '.'},
	{'#', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '#', '.', '.', '.'}}

func TestFromString(t *testing.T) {
	p, err := FromString(input)
	if err != nil {
		t.Errorf("FromString: %s", err)
		return
	}

	if err = sliceOfSlicesEqual(p.m, planMap); err != nil {
		t.Errorf("plan: %s, expected: %v, got %v", err, planMap, p.m)
	}
	if p.i != 6 || p.j != 4 {
		t.Errorf("expected guard to be at 6|4, got %d|%d", p.i, p.j)
	}
	if p.d != North {
		t.Errorf("expected direction to be %d, got %d", North, p.d)
	}
}

func TestAdvance(t *testing.T) {
	p, err := FromString("#.#\n...\n#.#")
	if err != nil {
		t.Errorf("FromString: %s", err)
		return
	}
	p.i, p.j, p.d = 1, 0, North

	tests := []struct {
		i, j int
		d    Direction

		expectI, expectJ int
		expectD          Direction
		done             bool
	}{
		{1, 0, North, 1, 1, East, false},
		{0, 1, East, 1, 1, South, false},
		{1, 2, South, 1, 1, West, false},
		{2, 1, West, 1, 1, North, false},
		{1, 1, East, 1, 2, East, false},
		{1, 2, East, 1, 3, East, true},
	}

	for i, test := range tests {
		p.i, p.j, p.d = test.i, test.j, test.d
		got, err := p.SimpleAdvance(true)
		if err != nil {
			t.Errorf("%d: error %s", i, err)
		}
		if p.i != test.expectI || p.j != test.expectJ || p.d != test.expectD {
			t.Errorf("expected guard to be at %d|%d facing %d, got %d|%d facing %d", test.expectI, test.expectJ, test.expectD, p.i, p.j, p.d)
		}
		if got != test.done {
			t.Errorf("expected guard to be done=%t, got %t", test.done, got)
		}
		if p.m[test.i][test.j] != Visited {
			t.Errorf("expected %d|%d to be visited, got %s", test.i, test.j, string(p.m[test.i][test.j]))
		}
	}
}

func TestCountVisits(t *testing.T) {
	p, err := FromString("###\n^.#\n...")
	if err != nil {
		t.Errorf("FromString: %s", err)
		return
	}

	for {
		done, err := p.SimpleAdvance(true)
		if err != nil {
			t.Errorf("error: %s", err)
			break
		}
		if done {
			break
		}
	}
	got := p.Count(Visited)
	expect := 3
	if got != expect {
		t.Errorf("expected %d visits, got %d", expect, got)
	}
}

func TestCountFakeObstacles(t *testing.T) {
	p, err := FromString("....\n.^.#\n#...\n..#.")
	if err != nil {
		t.Errorf("FromString: %s", err)
		return
	}

	for {
		done, err := p.SearchCirlces()
		if err != nil {
			t.Errorf("error: %s", err)
			break
		}
		if done {
			break
		}
	}
	expectMap := [][]rune{
		{'.', 'O', '.', '.'},
		{'.', 'X', '.', '#'},
		{'#', '.', '.', '.'},
		{'.', '.', '#', '.'},
	}
	if err = sliceOfSlicesEqual(p.m, expectMap); err != nil {
		t.Errorf("%s: expected %s, got %s", err, &Plan{m: expectMap}, p)
	}

	got := p.Count(Fake)
	expect := 1
	if got != expect {
		t.Errorf("expected %d visits, got %d", expect, got)
	}
}

func sliceOfSlicesEqual[T comparable](gotSliceSlice [][]T, expectSliceSlice [][]T) error {
	if len(gotSliceSlice) != len(expectSliceSlice) {
		return fmt.Errorf("expected slice of length %d, got %d", len(expectSliceSlice), len(gotSliceSlice))
	}
	for i, gotSlice := range gotSliceSlice {
		expectSlice := expectSliceSlice[i]
		if len(gotSlice) != len(expectSlice) {
			return fmt.Errorf("expected slice of length %d, got %d", len(expectSlice), len(gotSlice))
		}

		for j, got := range gotSlice {
			expect := expectSlice[j]
			if got != expect {
				return fmt.Errorf("at %d, %d: expected %v, got %v", i, j, expect, got)
			}
		}
	}
	return nil
}

func TestCopy(t *testing.T) {
	p, _ := FromString(input)
	c := p.Copy()
	c.m[p.i][p.j] = Guard
	if c.String() != input {
		t.Errorf("expected:\n%s\n\ngot:\n%s", input, c.String())
	}
}

// func TestOnCircle(t *testing.T) {
// 	g := NewGraph[*VisitData](4)
// 	g.AddEdge(0, 1, &VisitData{})
// 	g.AddEdge(1, 2, &VisitData{})
// 	g.AddEdge(2, 3, &VisitData{})
// 	g.AddEdge(3, 0, &VisitData{})
// 	got := OnCircle(g, 0)
// 	if !got {
// 		t.Errorf("circle not recognized")
// 	}
// }
