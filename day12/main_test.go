package main

import "testing"

func TestNewGarden(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"AABB\nCCDD"},
	}

	for _, test := range tests {
		g := NewGarden(test.input)
		got := g.String()
		if got != test.input {
			t.Errorf("expected %s, got %s", test.input, got)
		}
	}
}

func TestTotalPrice(t *testing.T) {
	tests := []struct {
		intput          string
		total, discount int
	}{
		{"A", 4, 4},
		{"AABB\nCCDD", 4 * 2 * 6, 4 * 2 * 4},
		{"AB\nBA", 4 * 1 * 4, 4 * 1 * 4},
		{"ABBB\nBCBB\nABBB", 4*1*4 + 8*14, 4*1*4 + 8*8},
	}

	for _, test := range tests {
		g := NewGarden(test.intput)
		total, discount := g.TotalPrice(0)
		if total != test.total {
			t.Errorf("expected total %d, got %d", test.total, total)
		}
		if discount != test.discount {
			t.Errorf("expected discount %d, got %d", test.discount, discount)
		}
	}
}

func TestRotateAntiClockwise(t *testing.T) {
	tests := []struct {
		start  Direction
		expect [2]int
	}{
		{North, West},
		{West, South},
		{South, East},
		{East, North},
	}

	for _, test := range tests {
		got := RotateAntiClockwise(test.start)
		if got != test.expect {
			t.Errorf("expected %v, got %v", test.expect, got)
		}
	}
}
