package main

import "testing"

func TestCountAfterNBlinks(t *testing.T) {
	tests := []struct {
		stones    []int
		maxBlinks int
		expect    int
	}{
		{[]int{0}, 1, 1},
		{[]int{10}, 1, 2},
		{[]int{1, 0}, 2, 3},
		{[]int{10}, 3, 3},
	}

	memory := make(StoneMemory)

	for _, test := range tests {
		got := CountAfterNBlinks(memory, test.stones, 0, test.maxBlinks)
		if got != test.expect {
			t.Errorf("exptected %d, got %d", test.expect, got)
		}
	}
}
