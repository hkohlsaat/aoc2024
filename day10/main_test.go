package main

import (
	"fmt"
	"testing"
)

func TestGraphFromString(t *testing.T) {
	input := `01234
94365
87430
06521`
	tests := []struct {
		id        int
		height    int
		neighbors map[int]bool
	}{
		{0, 0, map[int]bool{1: false, 5: false}},
		{3, 3, map[int]bool{2: false, 4: false, 8: false}},
	}

	heights, dimensions := ReadInput(input)
	g := GraphFromHeights(heights, dimensions)

	for _, test := range tests {
		node, ok := g.Nodes[test.id]
		if !ok {
			t.Errorf("node %d does not exist", test.id)
			continue
		}
		if node.Height != test.height {
			t.Errorf("node %d: expected height %d, got %d", test.id, test.height, node.Height)
			continue
		}
		if len(node.Neighbors) != len(test.neighbors) {
			t.Errorf("node %d: expected %d neighbors, got %d", test.id, len(test.neighbors), len(node.Neighbors))
			continue
		}
		for _, neighbor := range node.Neighbors {
			if _, ok := test.neighbors[neighbor.Id]; !ok {
				t.Errorf("got unexpected neighbor %v", neighbor)
			}
			test.neighbors[neighbor.Id] = true
		}
		for id, seen := range test.neighbors {
			if !seen {
				t.Errorf("did not see neighbor with id %d", id)
			}
		}
	}
}

func TestNodesWithHeight(t *testing.T) {
	input := `01234
94365
87430
06521`
	heights, dimensions := ReadInput(input)
	g := GraphFromHeights(heights, dimensions)

	expect := map[int]bool{
		PositionToNodeId([2]int{0, 0}, dimensions): false,
		PositionToNodeId([2]int{2, 4}, dimensions): false,
		PositionToNodeId([2]int{3, 0}, dimensions): false,
	}

	fmt.Printf("%v\n", heights)

	for id, node := range NodesWithHeight(g, 0) {
		if _, ok := expect[id]; !ok {
			t.Errorf("did not expect %d, %v", id, node)
			continue
		}
		expect[id] = true
	}

	for id, seen := range expect {
		if !seen {
			t.Errorf("did not get %d", id)
		}
	}
}

func TestScore(t *testing.T) {
	input := `01234
94365
87430
06521`
	heights, dimensions := ReadInput(input)
	g := GraphFromHeights(heights, dimensions)

	tests := []struct {
		id     int
		score  int
		rating int
	}{
		{PositionToNodeId([2]int{0, 0}, dimensions), 1, 1},
		{PositionToNodeId([2]int{2, 4}, dimensions), 1, 1},
		{PositionToNodeId([2]int{3, 1}, dimensions), 0, 0},
	}

	for _, test := range tests {
		gotScore, gotRating := TailheadMeasures(g, test.id)
		if gotScore != test.score {
			t.Errorf("expected score %d, got %d, for %v", test.score, gotScore, NodeIdToPosition(test.id, dimensions))
		}
		if gotRating != test.rating {
			t.Errorf("expected rating %d, got %d, for %v", test.rating, gotRating, NodeIdToPosition(test.id, dimensions))
		}
	}
}

func TestNodeIdPositionConversion(t *testing.T) {
	dimension := [2]int{10, 10}
	converted := PositionToNodeId(NodeIdToPosition(-10, dimension), dimension)
	if converted != -10 {
		t.Errorf("expected %d, got %d", -10, converted)
	}
}
