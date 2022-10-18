package graph

import (
	"fmt"

	"github.com/gammazero/deque"
)

type NodeNotFoundErr[K comparable] struct {
	Key K
}

func (e *NodeNotFoundErr[K]) Error() string {
	return fmt.Sprintf("node with a key %+v not found", e.Key)
}

// Graph contains all the nodes and edges
type Graph[K comparable, T any] struct {
	nodes map[K]T
	edges map[K][]K
}

func NewGraph[K comparable, T any]() Graph[K, T] {
	return Graph[K, T]{
		nodes: map[K]T{},
		edges: map[K][]K{},
	}
}

// SetNode sets a node with a value T to the key K
func (g *Graph[K, T]) SetNode(key K, value T) {
	g.nodes[key] = value
}

// GetEdges gets edges of node with key K
func (g *Graph[K, T]) GetEdges(key K) ([]K, error) {
	_, exists := g.nodes[key]
	if !exists {
		return nil, &NodeNotFoundErr[K]{Key: key}
	}

	edges, ok := g.edges[key]
	if ok {
		return edges, nil
	}

	g.edges[key] = []K{}

	return g.edges[key], nil
}

// AddEdge adds a directed edge between A and B (A -> B)
// If A already have B edge it will do nothing
func (g *Graph[K, T]) AddEdge(keyA K, keyB K) error {
	nodeAEdges, err := g.GetEdges(keyA)
	if err != nil {
		return err
	}

	if !edgeAlreadyPresent(nodeAEdges, keyB) {
		nodeAEdges = append(nodeAEdges, keyB)
		g.edges[keyA] = nodeAEdges
	}

	return nil
}

// AddUndirectedEdge adds an undirected edge between A and B (A <-> B)
// If A already have B edge or B already have A edge it will do nothing
func (g *Graph[K, T]) AddUndirectedEdge(keyA K, keyB K) error {
	nodeAEdges, err := g.GetEdges(keyA)
	if err != nil {
		return err
	}

	nodeBEdges, err := g.GetEdges(keyB)
	if err != nil {
		return err
	}

	if !edgeAlreadyPresent(nodeAEdges, keyB) {
		nodeAEdges = append(nodeAEdges, keyB)
		g.edges[keyA] = nodeAEdges
	}

	if !edgeAlreadyPresent(nodeBEdges, keyA) {
		nodeBEdges = append(nodeBEdges, keyA)
		g.edges[keyB] = nodeBEdges
	}

	return nil
}

// RemoveEdge removes an edge (another node with a key K) from the node with a key K
func (g *Graph[K, T]) RemoveEdge(key K, edge K) error {
	nodeEdges, err := g.GetEdges(key)
	if err != nil {
		return err
	}

	for i, e := range nodeEdges {
		if e == edge {
			nodeEdges = removeIndex(nodeEdges, i)
			g.edges[key] = nodeEdges
			return nil
		}
	}

	return fmt.Errorf("node key %+v doesn't have the edge %+v", key, edge)
}

// GetNode gets node T from key K
func (g *Graph[K, T]) GetNode(key K) (node T, err error) {
	node, ok := g.nodes[key]
	if ok {
		return node, nil
	}

	return node, &NodeNotFoundErr[K]{Key: key}
}

type path[K comparable] struct {
	nodeKey K
	prev    *path[K]
}

// Get shortest path between two node keys using breadth first search
func (g *Graph[K, T]) ShortestPath(start K, end K) ([]T, error) {
	queue := deque.New[path[K]]()
	queue.PushBack(path[K]{nodeKey: start, prev: nil})

	visited := map[K]bool{}
	for queue.Len() != 0 {
		p := queue.PopFront()
		if _, ok := visited[p.nodeKey]; ok {
			continue
		}

		visited[p.nodeKey] = true

		if p.nodeKey == end {
			return g.pathToArrayOfNodes(p), nil
		}

		nodeEdges, err := g.GetEdges(p.nodeKey)
		if err != nil {
			return nil, err
		}

		for _, edge := range nodeEdges {
			queue.PushBack(path[K]{nodeKey: edge, prev: &p})
		}
	}

	return []T{}, nil
}

// GetNodeKey invokes f on each node in the graph and returns a key as soon as f returns true.
// if f never returned true, returns (zeroValue for K), false
func (g *Graph[K, T]) GetNodeKey(f func(T) bool) (key K, ok bool) {
	for key, value := range g.nodes {
		if f(value) {
			return key, true
		}
	}

	var zeroValue K
	return zeroValue, false
}

// GetNodeKey invokes f on each node in the graph, if f returns true the current node key K
// that f is visiting will be added to the keys returned
func (g *Graph[K, T]) GetNodeKeys(f func(T) bool) (keys []K) {
	keys = []K{}

	for key, value := range g.nodes {
		if f(value) {
			keys = append(keys, key)
		}
	}

	return keys
}

func (g *Graph[K, T]) pathToArrayOfNodes(p path[K]) []T {
	nodes := []T{}
	for p.prev != nil {
		node, _ := g.GetNode(p.nodeKey)
		nodes = append(nodes, node)
		p = *p.prev
	}
	node, _ := g.GetNode(p.nodeKey)
	nodes = append(nodes, node) // last node

	// Reverse nodes
	reverseSlice(nodes)

	return nodes
}

func removeIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}

func reverseSlice[T any](input []T) {
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
}

func edgeAlreadyPresent[K comparable](nodeEdges []K, edgeKey K) bool {
	for _, e := range nodeEdges {
		if e == edgeKey {
			return true
		}
	}

	return false
}
