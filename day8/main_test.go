package main

import "testing"

const input = `............
........0...
.....0......
.......0....
....0.......
......A.....
............
............
........A...
.........A..
............
............`

func TestFrequenciesMapFromString(t *testing.T) {
	m := FrequenciesMapFromString(input)
	got := m.String()
	if got != input {
		t.Errorf("expected:\n%s\n\ngot:\n%s", input, got)
	}
}

func TestAntennaPositions(t *testing.T) {
	m := FrequenciesMapFromString(input)
	tests := []struct {
		frequency rune
		positions [][2]int
	}{
		{'A', [][2]int{{5, 6}, {8, 8}, {9, 9}}},
		{'0', [][2]int{{1, 8}, {2, 5}, {3, 7}, {4, 4}}},
	}

	antennaPositions := AntennaPositions(m)

	for _, test := range tests {
		positions, ok := antennaPositions[test.frequency]
		if !ok {
			t.Errorf("expected to have antenna with frequency %s, but was missing", string(test.frequency))
			continue
		}
		if len(positions) != len(test.positions) {
			t.Errorf("expected to get %d positions for %s-antennas, got %d", len(test.positions), string(test.frequency), len(positions))
			continue
		}
		for i := range positions {
			if positions[i][0] != test.positions[i][0] || positions[i][1] != test.positions[i][1] {
				t.Errorf("expected %v, got %v", test.positions[i], positions[i])
			}
		}
	}
}

func TestArithmetic(t *testing.T) {
	a := [2]int{2, 4}
	b := [2]int{6, 3}

	plusResult := Plus(a, b)
	if plusResult[0] != 8 || plusResult[1] != 7 {
		t.Errorf("expected [8 7], got %v", plusResult)
	}

	minusResult := Minus(a, b)
	if minusResult[0] != -4 || minusResult[1] != 1 {
		t.Errorf("expected [-4 1], got %v", minusResult)
	}
}

func TestAntinodePositionsDistinctDistances(t *testing.T) {
	tests := []struct {
		antennaPositions  map[rune][][2]int
		antinodePositions map[[2]int][]rune
	}{
		{map[rune][][2]int{'A': {{2, 3}, {3, 2}}, 'B': {{2, 2}, {4, 4}}}, map[[2]int][]rune{{0, 0}: {'B'}, {1, 4}: {'A'}, {4, 1}: {'A'}}},
	}

	for _, test := range tests {
		gotAntinodes := CollectAntinodePositions(test.antennaPositions, func(i1, i2 [2]int) func(yield func([2]int) bool) {
			return AntinodePositionsDistinctDistances(i1, i2, [2]int{5, 5})
		})
		if len(gotAntinodes) != len(test.antinodePositions) {
			t.Errorf("expected %d antinode positions, got %d: %v", len(test.antinodePositions), len(gotAntinodes), gotAntinodes)
			continue
		}
	positionsLoop:
		for position, expectedFrequencies := range test.antinodePositions {
			gotFrequencies, ok := gotAntinodes[position]
			if !ok {
				t.Errorf("expected to get antinode at %v, but got none", position)
				continue
			}
			if len(gotFrequencies) != len(expectedFrequencies) {
				t.Errorf("expected to get frequencies %v, got %v", expectedFrequencies, gotFrequencies)
				continue
			}
			for i := range expectedFrequencies {
				if gotFrequencies[i] != expectedFrequencies[i] {
					t.Errorf("expected to get frequencies %v, got %v", expectedFrequencies, gotFrequencies)
					continue positionsLoop
				}
			}
		}
	}
}

func TestAntinodePositionsExactlyInLine(t *testing.T) {
	tests := []struct {
		antennaPositions  map[rune][][2]int
		antinodePositions map[[2]int][]rune
	}{
		{map[rune][][2]int{'A': {{2, 3}, {3, 2}}, 'B': {{2, 2}, {4, 4}}}, map[[2]int][]rune{{0, 0}: {'B'}, {1, 1}: {'B'}, {2, 2}: {'B'}, {3, 3}: {'B'}, {4, 4}: {'B'}, {1, 4}: {'A'}, {2, 3}: {'A'}, {3, 2}: {'A'}, {4, 1}: {'A'}}},
	}

	for _, test := range tests {
		gotAntinodes := CollectAntinodePositions(test.antennaPositions, func(i1, i2 [2]int) func(yield func([2]int) bool) {
			return AntinodePositionsExactlyInLine(i1, i2, [2]int{5, 5})
		})
		if len(gotAntinodes) != len(test.antinodePositions) {
			t.Errorf("expected %d antinode positions, got %d: %v", len(test.antinodePositions), len(gotAntinodes), gotAntinodes)
			continue
		}
	positionsLoop:
		for position, expectedFrequencies := range test.antinodePositions {
			gotFrequencies, ok := gotAntinodes[position]
			if !ok {
				t.Errorf("expected to get antinode at %v, but got none", position)
				continue
			}
			if len(gotFrequencies) != len(expectedFrequencies) {
				t.Errorf("expected to get frequencies %v, got %v", expectedFrequencies, gotFrequencies)
				continue
			}
			for i := range expectedFrequencies {
				if gotFrequencies[i] != expectedFrequencies[i] {
					t.Errorf("expected to get frequencies %v, got %v", expectedFrequencies, gotFrequencies)
					continue positionsLoop
				}
			}
		}
	}
}

func TestEuclid(t *testing.T) {
	tests := []struct {
		a, b int
		gcd  int
	}{
		{1, 1, 1},
		{2, 0, 2},
		{3, 4, 1},
		{21, 14, 7},
	}

	for _, test := range tests {
		got := Euclid(test.a, test.b)
		if got != test.gcd {
			t.Errorf("expected %d, got %d", test.gcd, got)
		}
	}
}

func TestDiv(t *testing.T) {
	got := Div([2]int{4, 6}, 2)
	if got[0] != 2 || got[1] != 3 {
		t.Errorf("expected [2 3], got %v", got)
	}
}
