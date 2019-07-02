package top

import (
	"errors"
)

// ErrCyclicGraph indicates that a graph is cyclic. Only acyclic graphs can be
// sorted.
var ErrCyclicGraph = errors.New("cyclic graph")

// Graph to be sorted.
type Graph struct {
	nodes map[interface{}]*node
}

type mark byte

const (
	unmarked mark = iota
	temporary
	permanent
)

type node struct {
	of         interface{}
	mark       mark
	precursors []*node
}

// Link asserts that, when ordered, from should appear before to.
func (g *Graph) Link(from, to interface{}) {
	fromN := g.nodeFor(from)
	toN := g.nodeFor(to)
	toN.addPrecursor(fromN)
}

// Sort evaluates the graph to produce an ordering. If the graph cannot be sorted,
// an error is returned.
func (g *Graph) Sort() ([]interface{}, error) {
	sorted := make([]interface{}, 0, len(g.nodes))
	for _, n := range g.nodes {
		if n.mark != unmarked {
			continue
		}
		if err := n.visit(&sorted); err != nil {
			return nil, err
		}
	}
	return sorted, nil
}

func (g *Graph) nodeFor(x interface{}) *node {
	if g.nodes == nil {
		g.nodes = map[interface{}]*node{}
	}
	n := g.nodes[x]
	if n == nil {
		n = &node{of: x}
		g.nodes[x] = n
	}
	return n
}

func (n *node) addPrecursor(m *node) {
	for _, l := range n.precursors {
		if l == m {
			return
		}
	}
	n.precursors = append(n.precursors, m)
}

func (n *node) visit(sorted *[]interface{}) error {
	if n.mark == permanent {
		return nil
	}
	if n.mark == temporary {
		return ErrCyclicGraph
	}
	n.mark = temporary
	for _, n := range n.precursors {
		if err := n.visit(sorted); err != nil {
			return err
		}
	}
	n.mark = permanent
	*sorted = append(*sorted, n.of)
	return nil
}
