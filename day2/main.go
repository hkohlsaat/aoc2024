package main

import (
	"flag"
	"fmt"
	"os"
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

	reports, err := SplitReports(string(b))
	if err != nil {
		panic(fmt.Sprintf("extracting values from file: %s\n", err))
	}

	safeCount := 0
	for _, report := range reports {
		if isSave := CheckSafety(report); isSave {
			safeCount++
		}
	}

	fmt.Printf("Safe count: %d\n", safeCount)

	dampenedSafeCount := 0
	for _, report := range reports {
		if isSave := CheckDampenedSafety(report); isSave {
			dampenedSafeCount++
		}
	}
	fmt.Printf("With problem dampener: %d\n", dampenedSafeCount)
}

func SplitReports(s string) ([][]int, error) {
	lines := strings.Split(strings.TrimSpace(s), "\n")

	reports := make([][]int, 0, len(lines))

	for i, line := range lines {
		levels := strings.Split(line, " ")

		report := make([]int, 0, len(levels))

		for j, levelStr := range levels {
			level, err := strconv.Atoi(levelStr)
			if err != nil {
				return nil, fmt.Errorf("line %d, level %d: %w", i+1, j+1, err)
			}
			report = append(report, level)
		}

		reports = append(reports, report)
	}

	return reports, nil
}

func CheckSafety(report []int) bool {
	if len(report) < 2 {
		return true
	}

	direction := 1
	if report[0] > report[1] {
		direction = -1
	}

	for i := 0; i < len(report)-1; i++ {
		x, y := report[i], report[i+1]
		d := (y - x) * direction
		if d <= 0 || d > 3 {
			return false
		}
	}
	return true
}

func CheckDampenedSafety(report []int) bool {
	if len(report) < 3 {
		return true
	}

	differences := make([]int, len(report)-1)
	for i := range differences {
		differences[i] = report[i+1] - report[i]
	}

	zeroDifferenceIdx := -1
	for i, d := range differences {
		if d == 0 {
			if zeroDifferenceIdx != -1 {
				// We have seen one zero difference, already.
				// This is the second one. The dampener can not help here.
				return false
			}
			zeroDifferenceIdx = i
		}
	}

	if zeroDifferenceIdx != -1 {
		// We have seen a zero difference at this pair index.
		report = append(report[:zeroDifferenceIdx], report[zeroDifferenceIdx+1:]...)
		return CheckSafety(report)
	}

	if len(report) < 4 {
		return true
	}

	sgn := sign(sign(differences[0]) + sign(differences[1]) + sign(differences[2]))

	for i, difference := range differences {
		if sd := difference * sgn; sd >= 0 && sd < 4 {
			continue
		}

		c := append([]int(nil), report...)
		return CheckSafety(append(c[:i], c[i+1:]...)) || CheckSafety(append(report[:i+1], report[i+2:]...))
	}
	return true
}

func sign(i int) int {
	if i < 0 {
		return -1
	} else {
		return 1
	}
}
