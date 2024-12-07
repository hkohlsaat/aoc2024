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

	original, err := FromString(strings.TrimSpace(string(b)))
	if err != nil {
		panic(fmt.Sprintf("could not construct plan: %s", err))
	}

	plan := original.Copy()
	for {
		done, err := plan.SimpleAdvance(true)
		if err != nil {
			panic(fmt.Sprintf("error advancing: %s\n", err))
		}
		if done {
			break
		}
	}
	visits := plan.Count(Visited)

	fmt.Printf("Number of visits: %d\n", visits)

	plan = original.Copy()
	for {
		done, err := plan.SearchCirlces()
		if err != nil {
			panic(fmt.Sprintf("error advancing: %s\n", err))
		}
		if done {
			break
		}
	}
	fakeObstacles := plan.Count(Fake)
	if plan.m[original.i][original.j] == Fake {
		fakeObstacles -= 1
	}
	fmt.Printf("%s\n", plan)

	fmt.Printf("Number of fakes: %d\n", fakeObstacles)
}

type Plan struct {
	m              [][]rune
	i, j           int
	d              Direction
	startI, startJ int
}

type PlanGraph = Graph[*VisitData]

type VisitData struct {
	Visited bool
}

func NodeId(i, j int, d Direction, rowSize int) int {
	return 4*(i*rowSize+j) + int(d)
}

func (p *Plan) String() string {
	b := bytes.Buffer{}
	for i, row := range p.m {
		if i > 0 {
			b.WriteString("\n")
		}
		for _, block := range row {
			b.WriteString(string(block))
		}
	}
	return b.String()
}

func (p *Plan) Copy() *Plan {
	n := &Plan{
		m:      make([][]rune, len(p.m)),
		i:      p.i,
		j:      p.j,
		d:      p.d,
		startI: p.startI,
		startJ: p.startJ,
	}

	for i := range n.m {
		n.m[i] = make([]rune, len(p.m[i]))
		for j, block := range p.m[i] {
			if block == Visited {
				n.m[i][j] = Empty
			} else {
				n.m[i][j] = block
			}
		}
	}
	return n
}

const (
	Empty    rune = '.'
	Obstacle rune = '#'
	Guard    rune = '^'
	Visited  rune = 'X'
	Fake     rune = 'O'
)

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

func FromString(s string) (*Plan, error) {
	p := &Plan{}
	lines := strings.Split(s, "\n")
	p.m = make([][]rune, 0, len(lines))
	for i, line := range lines {
		row := make([]rune, 0, len(line))
		for j, ch := range line {
			if ch == Guard {
				p.i = i
				p.j = j
				p.d = North
				ch = Empty
			} else if ch != Empty && ch != Obstacle {
				return nil, fmt.Errorf("unknown block: %s", string(ch))
			}
			row = append(row, ch)
		}
		p.m = append(p.m, row)
	}
	return p, nil
}

func (p *Plan) SimpleAdvance(markVisitedOverNonEmpty bool) (done bool, err error) {
	if markVisitedOverNonEmpty || p.m[p.i][p.j] == Empty {
		p.m[p.i][p.j] = Visited
	}
	done, p.i, p.j, p.d, err = p.Plan(p.i, p.j, p.d)
	return done, err
}

func (p *Plan) Plan(i, j int, d Direction) (bool, int, int, Direction, error) {
	ni, nj := nextPosition(i, j, d)
	isObstacle, leaves := p.Inspect(ni, nj)

	if leaves {
		return true, ni, nj, d, nil
	}

	if isObstacle {
		for _, dd := range []Direction{1, 2, 3} {
			nd := (d + dd) % 4
			ni, nj = nextPosition(i, j, nd)
			isObstacle, _ := p.Inspect(ni, nj)
			if !isObstacle {
				return p.Plan(i, j, nd)
			}
		}
		return false, i, j, d, fmt.Errorf("could not move")
	}

	return false, ni, nj, d, nil
}

func (p *Plan) Inspect(i, j int) (isObstacle, leaves bool) {
	if i < 0 || i >= len(p.m) || j < 0 || j >= len(p.m[0]) {
		return false, true
	}

	if p.m[i][j] == Obstacle {
		return true, false
	}

	return false, false
}

func (p *Plan) SearchCirlces() (bool, error) {
	leave, ni, nj, _, err := p.Plan(p.i, p.j, p.d)
	if err != nil {
		return false, err
	}
	if leave || p.m[ni][nj] == Fake || p.m[ni][nj] == Visited {
		return p.SimpleAdvance(false)
	}

	if p.m[ni][nj] != Empty {
		return false, fmt.Errorf("expected empty block, got %q", string(p.m[ni][nj]))
	}

	p.m[ni][nj] = Obstacle
	circle, err := p.GoesInCircle(p.i, p.j, p.d)
	if err != nil {
		return false, err
	}
	if circle {
		p.m[ni][nj] = Fake
	} else {
		p.m[ni][nj] = Empty
	}

	return p.SimpleAdvance(false)
}

func (p *Plan) GoesInCircle(i, j int, d Direction) (leaves bool, err error) {
	started := NodeId(i, j, d, len(p.m))
	visited := make([]bool, 4*len(p.m)*len(p.m[0]))
	visited[started] = true
	for {
		leaves, i, j, d, err = p.Plan(i, j, d)
		if err != nil {
			return false, err
		}
		if leaves {
			return false, nil
		}
		n := NodeId(i, j, d, len(p.m))
		if visited[n] {
			return true, nil
		}
		visited[n] = true
	}
}

func (p *Plan) Count(target rune) int {
	count := 0
	for _, rows := range p.m {
		for _, block := range rows {
			if block == target {
				count += 1
			}
		}
	}
	return count
}

func nextPosition(i, j int, d Direction) (int, int) {
	switch d {
	case North:
		return i - 1, j
	case East:
		return i, j + 1
	case South:
		return i + 1, j
	case West:
		return i, j - 1
	}
	panic(fmt.Sprintf("unreachable code: unknown direction %d", d))
}
