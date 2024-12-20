package main

import (
	"bytes"
	"flag"
	"fmt"
	"iter"
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
	input := strings.TrimSpace(string(b))
	g := NewGarden(input)
	totalPrice, discountPrice := g.TotalPrice(0)
	fmt.Printf("Total price: %d\n", totalPrice)
	fmt.Printf("Discount price: %d\n", discountPrice)
}

type Garden struct {
	Blocks  []rune
	Visited []bool

	Dimensions [2]int
}

func NewGarden(input string) *Garden {
	lines := strings.Split(input, "\n")
	if len(lines) == 0 {
		panic("lines are empty")
	}

	dimensions := [2]int{
		len(lines),
		len(lines[0]),
	}

	blocksLength := dimensions[0] * dimensions[1]
	blocks := make([]rune, 0, blocksLength)

	for _, line := range lines {
		for _, r := range line {
			blocks = append(blocks, r)
		}
	}

	if len(blocks) != blocksLength {
		panic(fmt.Sprintf("did not read blocks according to dimensions, expected %d, got %d",
			blocksLength, len(blocks)))
	}

	return &Garden{
		Blocks:  blocks,
		Visited: make([]bool, blocksLength),

		Dimensions: dimensions,
	}
}

func (g *Garden) PositionToId(position [2]int) int {
	return position[0]*g.Dimensions[1] + position[1]
}

func (g *Garden) At(position [2]int) rune {
	if !g.ValidPosition(position) {
		panic(fmt.Sprintf("position %v is not valid", position))
	}
	return g.Blocks[g.PositionToId(position)]
}

func (g *Garden) IdToPosition(id int) [2]int {
	return [2]int{
		id / g.Dimensions[1],
		id % g.Dimensions[1],
	}
}

func (g *Garden) ValidPosition(position [2]int) bool {
	return position[0] >= 0 && position[0] < g.Dimensions[0] && position[1] >= 0 && position[1] < g.Dimensions[1]
}

func (g *Garden) String() string {
	b := &bytes.Buffer{}

	for i := range g.Dimensions[0] {
		if i > 0 {
			b.WriteRune('\n')
		}
		for j := range g.Dimensions[1] {
			id := g.PositionToId([2]int{i, j})
			b.WriteRune(g.Blocks[id])
		}
	}

	return b.String()
}

func (g *Garden) TotalPrice(start int) (totalPrice, discountPrice int) {
	if g.Visited[start] {
		return 0, 0
	}

	var (
		area      = 0
		perimeter = 0
		sides     = 0
	)

	otherAreaBlocks := make([]int, 0, len(g.Blocks))
	ownAreaBlocks := append(make([]int, 0, len(g.Blocks)), start)

	for len(ownAreaBlocks) > 0 {
		id := ownAreaBlocks[0]
		ownAreaBlocks = ownAreaBlocks[1:]

		if g.Visited[id] {
			continue
		}
		g.Visited[id] = true

		area += 1

		position := g.IdToPosition(id)
		for direction, neighboringPosition := range SimpleNeighboringPositions(position) {
			neighboringId := g.PositionToId(neighboringPosition)

			if !g.ValidPosition(neighboringPosition) {
				perimeter += 1
				if g.HasSide(id, RotateAntiClockwise(direction)) {
					sides += 1
				} else {
					leftNeighborPosition := Plus(position, RotateAntiClockwise(direction))
					leftNeighborId := g.PositionToId(leftNeighborPosition)
					if !g.HasSide(leftNeighborId, direction) {
						sides += 1
					}
				}

			} else if g.Blocks[neighboringId] != g.Blocks[id] {
				perimeter += 1
				if g.HasSide(id, RotateAntiClockwise(direction)) {
					sides += 1
				} else {
					leftNeighborPosition := Plus(position, RotateAntiClockwise(direction))
					leftNeighborId := g.PositionToId(leftNeighborPosition)
					if !g.HasSide(leftNeighborId, direction) {
						sides += 1
					}
				}
				otherAreaBlocks = append(otherAreaBlocks, neighboringId)
			} else if !g.Visited[neighboringId] {
				ownAreaBlocks = append(ownAreaBlocks, neighboringId)
			}
		}
	}

	totalPrice = area * perimeter
	discountPrice = area * sides

	for _, block := range otherAreaBlocks {
		t, d := g.TotalPrice(block)
		totalPrice += t
		discountPrice += d
	}

	return totalPrice, discountPrice
}

func (g *Garden) HasSide(id int, direction Direction) bool {
	position := g.IdToPosition(id)
	neighborPosition := Plus(position, direction)
	if !g.ValidPosition(neighborPosition) {
		return true
	}
	neighborId := g.PositionToId(neighborPosition)
	return g.Blocks[neighborId] != g.Blocks[id]
}

type Direction = [2]int

var (
	North Direction = [2]int{-1, 0}
	East  Direction = [2]int{0, 1}
	South Direction = [2]int{1, 0}
	West  Direction = [2]int{0, -1}
)

func SimpleNeighboringPositions(position [2]int) iter.Seq2[Direction, [2]int] {
	return func(yield func([2]int, [2]int) bool) {
		for _, direction := range []Direction{North, East, South, West} {
			if !yield(direction, Plus(position, direction)) {
				return
			}
		}
	}
}

func Plus(v0, v1 [2]int) [2]int {
	return [2]int{
		v0[0] + v1[0],
		v0[1] + v1[1],
	}
}

func RotateAntiClockwise(v [2]int) [2]int {
	return [2]int{
		-v[1],
		v[0],
	}
}
