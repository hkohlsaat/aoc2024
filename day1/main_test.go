package main

import (
	"testing"
)

func TestListDifference(t *testing.T) {
	tests := []struct {
		Input
		expect int
	}{
		{Input{[]int{1}, []int{2}}, 1},
		{Input{[]int{3, 4, 2}, []int{4, 3, 5}}, 3},
		{Input{[]int{3, 4, 2, 1, 3, 3}, []int{4, 3, 5, 3, 9, 3}}, 11},
	}

	for _, input := range tests {
		got, err := ListDifference(input.left, input.right)
		if err != nil {
			t.Errorf("returned error: %s", err)
			continue
		}
		if got != input.expect {
			t.Errorf("expected %d, got %d", input.expect, got)
		}
	}
}

func TestSplitLists(t *testing.T) {
	tests := []struct {
		raw    string
		expect Input
	}{
		{"1   3", Input{[]int{1}, []int{3}}},
		{`3   4
4   3
2   5
1   3
3   9
3   3
`, Input{[]int{3, 4, 2, 1, 3, 3}, []int{4, 3, 5, 3, 9, 3}},
		},
	}

	for _, test := range tests {
		got, err := SplitLists(test.raw)
		if err != nil {
			t.Errorf("returned error: %s", err)
			continue
		}
		if !checkListsEqual(got.left, test.expect.left) {
			t.Errorf("left list: expected %v, got %v", test.expect.left, got.left)
		}
		if !checkListsEqual(got.right, test.expect.right) {
			t.Errorf("right list: expected %v, got %v", test.expect.right, got.right)
		}
	}
}

func checkListsEqual(l1, l2 []int) bool {
	if len(l1) != len(l2) {
		return false
	}
	for i := range l1 {
		if l1[i] != l2[i] {
			return false
		}
	}
	return true
}