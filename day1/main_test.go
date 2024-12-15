package main

import (
	"fmt"
	"os"
	"testing"
)

var input = func() string {
	filename := "input_test.txt"
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("could not read file %s: %s\n", filename, err))
	}
	return string(b)
}()

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
		{input, Input{[]int{3, 4, 2, 1, 3, 3}, []int{4, 3, 5, 3, 9, 3}}},
	}

	for _, test := range tests {
		got, err := SplitLists(test.raw)
		if err != nil {
			t.Errorf("returned error: %s", err)
			continue
		}
		if !areListsEqual(got.left, test.expect.left) {
			t.Errorf("left list: expected %v, got %v", test.expect.left, got.left)
		}
		if !areListsEqual(got.right, test.expect.right) {
			t.Errorf("right list: expected %v, got %v", test.expect.right, got.right)
		}
	}
}

func areListsEqual(l1, l2 []int) bool {
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

func TestSimilarityScore(t *testing.T) {
	tests := []struct {
		Input
		expect int
	}{
		{Input{[]int{1}, []int{1}}, 1},
		{Input{[]int{3, 4, 2, 1, 3, 3}, []int{4, 3, 5, 3, 9, 3}}, 31},
	}

	for _, input := range tests {
		score, err := SimilarityScore(input.left, input.right)
		if err != nil {
			t.Errorf("returned error: %s", err)
			continue
		}
		if score != input.expect {
			t.Errorf("expected %d, got %d", input.expect, score)
		}
	}
}
