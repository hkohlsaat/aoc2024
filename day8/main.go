package main

import (
	"bytes"
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

	m := FrequenciesMapFromString(string(b))
	antennaPositions := AntennaPositions(m)
	antinodePoisions := CollectAntinodePositions(antennaPositions, func(i1, i2 [2]int) func(yield func([2]int) bool) {
		return AntinodePositionsDistinctDistances(i1, i2, [2]int{len(m), len(m[0])})
	})
	fmt.Printf("Distinct antinode positions at distinct distances: %d\n", len(antinodePoisions))

	antinodePoisions = CollectAntinodePositions(antennaPositions, func(i1, i2 [2]int) func(yield func([2]int) bool) {
		return AntinodePositionsExactlyInLine(i1, i2, [2]int{len(m), len(m[0])})
	})
	fmt.Printf("Distinct antinode positions exactly in line: %d\n", len(antinodePoisions))
}

type Map [][]rune

const Empty = '.'

func FrequenciesMapFromString(s string) Map {
	lines := strings.Split(s, "\n")
	frequenciesMap := make([][]rune, len(lines))
	for i, line := range lines {
		frequenciesMap[i] = make([]rune, len(line))
		for j, block := range line {
			frequenciesMap[i][j] = block
		}
	}
	return frequenciesMap
}

func (m Map) String() string {
	b := &bytes.Buffer{}
	for i := range m {
		if i > 0 {
			b.WriteRune('\n')
		}
		for j := range m[i] {
			b.WriteRune(m[i][j])
		}
	}
	return b.String()
}

func AntennaPositions(m Map) map[rune][][2]int {
	antennaPositions := map[rune][][2]int{}
	for i := range m {
		for j := range m[i] {
			if m[i][j] == Empty {
				continue
			}
			antennaPositions[m[i][j]] = append(antennaPositions[m[i][j]], [2]int{i, j})
		}
	}
	return antennaPositions
}

func CollectAntinodePositions(antennaPositions map[rune][][2]int, positionAntinodes func([2]int, [2]int) func(yield func([2]int) bool)) map[[2]int][]rune {
	antinodePositions := map[[2]int][]rune{}
	for frequency, positions := range antennaPositions {
		for a, positionA := range positions {
			for _, positionB := range positions[a+1:] {
				for antinodePosition := range positionAntinodes(positionA, positionB) {
					antinodePositions[antinodePosition] = append(antinodePositions[antinodePosition], frequency)
				}
			}
		}
	}
	return antinodePositions
}

func AntinodePositionsDistinctDistances(antennaPosition1, antennaPosition2, dimensions [2]int) func(yield func([2]int) bool) {
	return func(yield func([2]int) bool) {
		difference := Minus(antennaPosition2, antennaPosition1)

		if antinodePosition := Minus(antennaPosition1, difference); ValidPosition(antinodePosition, dimensions) {
			if !yield(antinodePosition) {
				return
			}
		}

		if antinodePosition := Plus(antennaPosition2, difference); ValidPosition(antinodePosition, dimensions) {
			if !yield(antinodePosition) {
				return
			}
		}
	}
}

func AntinodePositionsExactlyInLine(antennaPosition1, antennaPosition2, dimensions [2]int) func(yield func([2]int) bool) {
	return func(yield func([2]int) bool) {
		difference := Minus(antennaPosition2, antennaPosition1)
		gcd := Euclid(Abs(difference[0]), Abs(difference[1]))
		difference = Div(difference, gcd)

		for i := 0; true; i++ {
			antinodePosition := Plus(antennaPosition1, Mul(difference, i))
			if !ValidPosition(antinodePosition, dimensions) {
				break
			}
			if !yield(antinodePosition) {
				return
			}
		}

		for i := -1; true; i-- {
			antinodePosition := Plus(antennaPosition1, Mul(difference, i))
			if !ValidPosition(antinodePosition, dimensions) {
				break
			}
			if !yield(antinodePosition) {
				return
			}
		}
	}
}

func Euclid(a, b int) int {
	if a < b {
		return Euclid(b, a)
	}
	if b == 0 {
		return a
	}
	return Euclid(b, a%b)
}

func Minus(vector1, vector2 [2]int) [2]int {
	return [2]int{
		vector1[0] - vector2[0],
		vector1[1] - vector2[1],
	}
}

func Plus(vector1, vector2 [2]int) [2]int {
	return [2]int{
		vector1[0] + vector2[0],
		vector1[1] + vector2[1],
	}
}

func Div(vector [2]int, divisor int) [2]int {
	return [2]int{
		vector[0] / divisor,
		vector[1] / divisor,
	}
}

func Mul(vector [2]int, factor int) [2]int {
	return [2]int{
		vector[0] * factor,
		vector[1] * factor,
	}
}

func Abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func ValidPosition(position, dimensions [2]int) bool {
	return position[0] >= 0 && position[0] < dimensions[0] && position[1] >= 0 && position[1] < dimensions[1]
}
