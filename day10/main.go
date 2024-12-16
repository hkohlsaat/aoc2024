package main

import (
	"bytes"
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
	heights, dimensions := ReadInput(string(b))
	g := GraphFromHeights(heights, dimensions)
	scoresSum, ratingsSum := SumTrailheadScores(g)
	fmt.Printf("Sum of trailhead scores: %d\n", scoresSum)
	fmt.Printf("Sum of trailhead ratings: %d\n", ratingsSum)
}

func ReadInput(input string) (heights []int, dimensions [2]int) {
	if len(input) == 0 {
		return nil, [2]int{0, 0}
	}
	lines := strings.Split(input, "\n")
	dimensions = [2]int{len(lines), len(lines[0])}
	for _, line := range lines {
		for _, r := range line {
			if r == '.' {
				heights = append(heights, -1)
				continue
			}

			height, err := strconv.Atoi(string(r))
			if err != nil {
				panic(err)
			}
			heights = append(heights, height)
		}
	}
	if len(heights) != dimensions[0]*dimensions[1] {
		panic(fmt.Sprintf("Length of heights does not match computed dimensions %v.", dimensions))
	}
	return heights, dimensions
}

type Graph struct {
	Nodes map[int]*Node
}

type Node struct {
	NodeInfo
	Id        int
	Neighbors []*Node
}

type NodeInfo struct {
	Height int
	Trails map[int][][]int
}

func GraphFromHeights(heights []int, dimensions [2]int) *Graph {
	g := &Graph{Nodes: make(map[int]*Node)}
	for id, height := range heights {
		g.Nodes[id] = &Node{
			Id:        id,
			NodeInfo:  NodeInfo{Height: height, Trails: map[int][][]int{}},
			Neighbors: make([]*Node, 0, 4),
		}
	}

	for id, node := range g.Nodes {
		position := NodeIdToPosition(id, dimensions)
		for _, neighborCandidate := range [][2]int{
			{position[0] - 1, position[1]},
			{position[0], position[1] + 1},
			{position[0] + 1, position[1]},
			{position[0], position[1] - 1},
		} {
			if !ValidPosition(neighborCandidate, dimensions) {
				continue
			}
			neighborId := PositionToNodeId(neighborCandidate, dimensions)
			node.Neighbors = append(node.Neighbors, g.Nodes[neighborId])
		}
		// for _, d := range []int{1, -dimensions[1], dimensions[1], -1} {
		// 	neighborId := id + d
		// 	if !ValidPosition(NodeIdToPosition(neighborId, dimensions), dimensions) {
		// 		continue
		// 	}
		// 	node.Neighbors = append(node.Neighbors, g.Nodes[neighborId])
		// }
	}
	return g
}

func PositionToNodeId(position, dimensions [2]int) int {
	return position[0]*dimensions[1] + position[1]
}

func NodeIdToPosition(id int, dimensions [2]int) [2]int {
	return [2]int{
		id / dimensions[1],
		id % dimensions[1],
	}
}

func ValidPosition(position, dimensions [2]int) bool {
	return position[0] >= 0 && position[0] < dimensions[0] && position[1] >= 0 && position[1] < dimensions[1]
}

func NodesWithHeight(g *Graph, targetHeight int) func(yield func(id int, node *Node) bool) {
	return func(yield func(id int, node *Node) bool) {
		for id := range len(g.Nodes) {
			node := g.Nodes[id]
			if node.Height == targetHeight {
				if !yield(id, node) {
					return
				}
			}
		}
	}
}

func TailheadMeasures(g *Graph, trailheadId int) (score, rating int) {
	for trail := range func(yield func([]int) bool) {
		BuildTrails(g, make([]int, 0, 10), trailheadId, func(from, to *Node) bool { return to.Height-from.Height == 1 }, yield)
	} {
		if trail[0] != trailheadId {
			panic("trail does not begin at correct starting position")
		}
		endOfTrail := trail[len(trail)-1]
		g.Nodes[trailheadId].Trails[endOfTrail] = append(g.Nodes[trailheadId].Trails[endOfTrail], trail)
		rating++
	}
	return len(g.Nodes[trailheadId].Trails), rating
}

func BuildTrails(g *Graph, trail []int, next int, valid func(from, to *Node) bool, yield func([]int) bool) bool {
	trail = append(trail, next)
	if len(trail) == 10 {
		return yield(append([]int(nil), trail...))
	}
	node := g.Nodes[next]
	for _, neighbor := range node.Neighbors {
		if valid(node, neighbor) {
			if !BuildTrails(g, trail, neighbor.Id, valid, yield) {
				return false
			}
		}
	}
	return true
}

func SumTrailheadScores(g *Graph) (scoresSum, ratingsSum int) {

	for id, _ := range NodesWithHeight(g, 0) {
		score, rating := TailheadMeasures(g, id)
		scoresSum += score
		ratingsSum += rating
	}
	return scoresSum, ratingsSum
}

func TrailString(trail []int, dimensions [2]int) string {
	runes := make([][]rune, dimensions[0])
	for i := range runes {
		runes[i] = make([]rune, dimensions[1])
	}

	for i, id := range trail {
		position := NodeIdToPosition(id, dimensions)
		runes[position[0]][position[1]] = '0' + rune(i)
	}
	b := &bytes.Buffer{}
	for i := range dimensions[1] + 2 {
		if i == 0 || i == dimensions[1]+1 {
			b.WriteRune('+')
		} else {
			b.WriteRune('-')
		}
	}
	b.WriteRune('\n')
	for i, line := range runes {
		if i == 0 {
			b.WriteRune('|')
		} else {
			b.WriteString("|\n|")
		}
		for _, r := range line {
			if r == 0 {
				b.WriteRune(' ')
			} else {
				b.WriteRune(r)
			}
		}
	}
	b.WriteString("|\n")
	for i := range dimensions[1] + 2 {
		if i == 0 || i == dimensions[1]+1 {
			b.WriteRune('+')
		} else {
			b.WriteRune('-')
		}
	}
	return b.String()
}
