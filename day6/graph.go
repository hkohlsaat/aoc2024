package main

type Graph[T any] struct {
	Nodes map[int]*Node[T]
}

func NewGraph[T any](nodes int) *Graph[T] {
	n := map[int]*Node[T]{}
	for i := 0; i < nodes; i++ {
		n[i] = &Node[T]{map[int]T{}}
	}
	return &Graph[T]{n}
}

func (g *Graph[T]) AddEdge(from, to int, data T) {
	f := g.Nodes[from]
	f.OutNeighbors[to] = data
}

type Node[T any] struct {
	OutNeighbors map[int]T
}
