package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var filename = flag.String("input", "input.txt", "input for this assignment")

func main() {
	flag.Parse()

	b, err := os.ReadFile(*filename)
	if err != nil {
		panic(fmt.Sprintf("could not read file %s: %s\n", *filename, err))
	}

	input := ByteSclices(strings.TrimSpace(string(b)))

	count := CountWords("XMAS", input)
	fmt.Printf("simple XMAS counter: %d\n", count)

	xmasCount := CountXMas(input)
	fmt.Printf("cross XMAS counter: %d\n", xmasCount)
}

func ByteSclices(input string) [][]byte {
	lines := strings.Split(input, "\n")
	for i := 1; i < len(lines); i++ {
		if len(lines[i]) != len(lines[i-1]) {
			panic(fmt.Sprintf("could not transform into [][]byte: lines %d and %d have different lengths", i-1, i))
		}
	}
	out := make([][]byte, len(lines))
	for i, line := range lines {
		out[i] = []byte(line)
	}
	return out
}

func CountWords(w string, input [][]byte) int {
	counter := 0
	locs := FindLetters(w[0], input)
	for _, loc := range locs {
		for _, dx := range []int{1, 0, -1} {
			for _, dy := range []int{1, 0, -1} {
				if dx == 0 && dy == 0 {
					continue
				}
				if CompleteWord([]byte(w), input, loc[0], loc[1], dx, dy) {
					counter += 1
				}
			}
		}
	}
	return counter
}

func FindLetters(r byte, input [][]byte) [][]int {
	locs := [][]int{}
	for i, line := range input {
		for j, b := range line {
			if b == r {
				locs = append(locs, []int{i, j})
			}
		}
	}
	return locs
}

func CompleteWord(w []byte, input [][]byte, x, y, dx, dy int) bool {
	for _, b := range w {
		if x < 0 || x >= len(input) || y < 0 || y >= len(input[0]) {
			return false
		}
		if input[x][y] != b {
			return false
		}
		x += dx
		y += dy
	}
	return true
}

func CountXMas(input [][]byte) int {
	counter := 0
	locs := FindLetters('A', input)
	for _, loc := range locs {
		if CompleteXMas(loc[0], loc[1], input) {
			counter += 1
		}
	}
	return counter
}

func CompleteXMas(x, y int, input [][]byte) bool {
	tlbr := (CheckByte(x-1, y-1, input, 'M') &&
		CheckByte(x+1, y+1, input, 'S')) ||
		(CheckByte(x-1, y-1, input, 'S') &&
			CheckByte(x+1, y+1, input, 'M'))
	rtbl := (CheckByte(x-1, y+1, input, 'M') &&
		CheckByte(x+1, y-1, input, 'S')) ||
		(CheckByte(x-1, y+1, input, 'S') &&
			CheckByte(x+1, y-1, input, 'M'))
	return tlbr && rtbl
}

func CheckByte(x, y int, input [][]byte, b byte) bool {
	if x < 0 || x >= len(input) || y < 0 || y >= len(input[0]) {
		return false
	}
	return input[x][y] == b
}
