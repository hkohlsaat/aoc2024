package main

import (
	"fmt"
	"os"
	"testing"
)

func TestCountWords(t *testing.T) {
	tests := []struct {
		w      string
		input  [][]byte
		expect int
	}{
		{"XMAS", ByteSclices("X###\n#M##\n##A#\n###S"), 1},
		{"XMAS", ByteSclices(`MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX`), 18},
	}

	for _, test := range tests {
		got := CountWords(test.w, test.input)
		if got != test.expect {
			t.Errorf("expected %d, got %d", test.expect, got)
		}
	}
}

func TestFindLetter(t *testing.T) {
	tests := []struct {
		letter byte
		input  [][]byte
		expect [][]int
	}{
		{'X', ByteSclices("X.\nYZ"), [][]int{{0, 0}}},
	}

	for _, test := range tests {
		locs := FindLetters(test.letter, test.input)
		if len(locs) != len(test.expect) {
			t.Errorf("expected %d locations, got %d", len(test.expect), len(locs))
			continue
		}
		for i, loc := range locs {
			if len(loc) != 2 {
				t.Errorf("expected location of length 2, got %d", len(loc))
			}
			if loc[0] != test.expect[i][0] || loc[1] != test.expect[i][1] {
				t.Errorf("location at index %d: expected %v, got %v", i, test.expect[i], loc)
			}
		}
	}
}

func TestCompleteWord(t *testing.T) {
	tests := []struct {
		w            []byte
		input        [][]byte
		x, y, dx, dy int
		expect       bool
	}{
		{[]byte("XMAS"), ByteSclices("X###\n#M##\n##A#\n###S"), 0, 0, 1, 1, true},
	}

	for _, test := range tests {
		completed := CompleteWord(test.w, test.input, test.x, test.y, test.dx, test.dy)
		if completed != test.expect {
			t.Errorf("expected %t, got %t", test.expect, completed)
		}
	}
}

func TestByteSlices(t *testing.T) {
	s := ByteSclices("X###\n#M##\n##A#\n###S")
	expect := [][]byte{[]byte("X###"), []byte("#M##"), []byte("##A#"), []byte("###S")}

	if len(s) != len(expect) {
		t.Errorf("expected height %d, got %d", len(expect), len(s))
		return
	}
	if len(s[0]) != len(expect[0]) {
		t.Errorf("expected width %d, got %d", len(expect[0]), len(s[0]))
		return
	}

	for i, line := range s {
		for j, b := range line {
			if b != expect[i][j] {
				t.Errorf("at %d, %d: expected %v, %v", i, j, expect[i][j], b)
			}
		}
	}
}

func TestCheckByte(t *testing.T) {
	tests := []struct {
		x, y   int
		input  [][]byte
		b      byte
		expect bool
	}{
		{0, 0, [][]byte{{'X'}}, 'X', true},
		{-1, 0, [][]byte{{'X'}}, 'X', false},
		{0, -1, [][]byte{{'X'}}, 'X', false},
		{0, 3, [][]byte{{'X'}}, 'X', false},
		{1, 1, [][]byte{{'X'}}, 'X', false},
	}

	for _, test := range tests {
		got := CheckByte(test.x, test.y, test.input, test.b)
		if got != test.expect {
			t.Errorf("expected %t, got %t", test.expect, got)
		}
	}
}

func TestCompleteXMas(t *testing.T) {
	tests := []struct {
		x, y   int
		input  [][]byte
		expect bool
	}{
		{2, 2, ByteSclices("X###\n#M#M\n##A#\n#S#S"), true},
	}

	for _, test := range tests {
		completed := CompleteXMas(test.x, test.y, test.input)
		if completed != test.expect {
			t.Errorf("expected %t, got %t", test.expect, completed)
		}
	}
}

func TestCountXMas(t *testing.T) {
	filename := "input_test.txt"
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("could not read file %s: %s\n", filename, err))
	}
	tests := []struct {
		input  [][]byte
		expect int
	}{
		{ByteSclices("X###\n#M#M\n##A#\n#S#S"), 1},
		{ByteSclices(string(b)), 9},
	}

	for _, test := range tests {
		got := CountXMas(test.input)
		if got != test.expect {
			t.Errorf("expected %d, got %d", test.expect, got)
		}
	}
}
