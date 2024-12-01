package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var filename = flag.String("input", "input.txt", "input for this assignment")

func main() {
	flag.Parse()

	b, err := os.ReadFile(*filename)
	if err != nil {
		panic(fmt.Sprintf("could not read file %s: %s\n", *filename, err))
	}

	s := string(b)
	input, err := SplitLists(s)
	if err != nil {
		panic(fmt.Sprintf("extracting values from file: %s\n", err))
	}

	totalDifference, err := ListDifference(input.left, input.right)
	if err != nil {
		panic(fmt.Sprintf("calculating total difference: %s\n", err))
	}

	fmt.Println(totalDifference)
}

func SplitLists(s string) (Input, error) {
	lines := strings.Split(strings.TrimSpace(s), "\n")

	input := Input{
		left:  make([]int, 0, len(lines)),
		right: make([]int, 0, len(lines)),
	}
	for i, line := range lines {
		entries := strings.Split(line, " ")
		leftStr := entries[0]
		rightStr := entries[len(entries)-1]
		left, err := strconv.Atoi(leftStr)
		if err != nil {
			return Input{}, fmt.Errorf("line %d (left column): %w", i+1, err)
		}

		right, err := strconv.Atoi(rightStr)
		if err != nil {
			return Input{}, fmt.Errorf("line %d (right column): %w", i+1, err)
		}

		input.left = append(input.left, left)
		input.right = append(input.right, right)
	}
	return input, nil
}

type Input struct {
	left, right []int
}

func ListDifference(left, right []int) (int, error) {
	if len(left) != len(right) {
		return 0, fmt.Errorf("lists' lengths aren't equal: left is %d, right is %d", len(left), len(right))
	}
	slices.SortFunc(left, func(l, r int) int { return r - l })
	slices.SortFunc(right, func(l, r int) int { return r - l })

	total := 0

	for i := range left {
		l := left[i]
		r := right[i]

		d := r - l

		// take the abs
		if d < 0 {
			d = -d
		}

		total += d
	}

	return total, nil
}
