package graph_test

import (
	"testing"

	"github.com/joetifa2003/graph-go"
	"github.com/stretchr/testify/assert"
)

func TestGetNode(t *testing.T) {
	assert := assert.New(t)

	g := graph.NewGraph[int, string]()
	g.SetNode(1, "Node 1")

	_, err := g.GetNode(3) // non exiting node
	assert.Error(err)
}

func TestGraphEdges(t *testing.T) {
	assert := assert.New(t)

	g := graph.NewGraph[int, string]()
	g.SetNode(1, "Node 1")
	g.SetNode(2, "Node 2")
	err := g.AddEdge(1, 2)
	assert.Nil(err)
	err = g.AddEdge(2, 1)
	assert.Nil(err)
	err = g.AddEdge(2, 1) // duplicate edge
	assert.Nil(err)
	err = g.AddEdge(69, 145) // none existing node
	assert.Equal(err.Error(), (&graph.NodeNotFoundErr[int]{Key: 69}).Error())
	assert.Error(err)

	node1, err := g.GetNode(1)
	assert.Nil(err)
	node1Edges, err := g.GetEdges(1)
	assert.Nil(err)
	node2, err := g.GetNode(2)
	assert.Nil(err)
	node2Edges, err := g.GetEdges(2)
	assert.Nil(err)

	assert.Equal("Node 1", node1)
	assert.Equal("Node 2", node2)
	assert.Equal(1, len(node1Edges))
	assert.Equal(1, len(node2Edges))

	err = g.RemoveEdge(1, 2)
	assert.Nil(err)
	node1Edges, err = g.GetEdges(1)
	assert.Nil(err)
	assert.Equal(0, len(node1Edges))

	err = g.RemoveEdge(1, 2) // already removed edge
	assert.Error(err)

	err = g.RemoveEdge(69, 145) // none existing node
	assert.Error(err)
}

func TestShortestPath(t *testing.T) {
	assert := assert.New(t)

	g := graph.NewGraph[int, string]()
	g.SetNode(1, "Node 1")
	g.SetNode(2, "Node 2")
	g.SetNode(3, "Node 3")
	g.SetNode(4, "Node 4")
	g.SetNode(5, "Node 5")
	g.SetNode(6, "Node 6")
	err := g.AddEdge(1, 2)
	assert.Nil(err)
	err = g.AddEdge(2, 1)
	assert.Nil(err)
	err = g.AddEdge(2, 3)
	assert.Nil(err)
	err = g.AddEdge(3, 2)
	assert.Nil(err)
	err = g.AddEdge(3, 4)
	assert.Nil(err)
	err = g.AddEdge(3, 1)
	assert.Nil(err)
	err = g.AddEdge(4, 3)
	assert.Nil(err)
	err = g.AddEdge(4, 5)
	assert.Nil(err)

	path1, err := g.ShortestPath(1, 4)
	assert.Nil(err)
	path2, err := g.ShortestPath(4, 6)
	assert.Nil(err)

	assert.Equal(4, len(path1))
	assert.Equal(0, len(path2)) // no possible path

	_, err = g.ShortestPath(69, 145) // non existing nodes
	assert.Error(err)
}

func TestUndirectedEdge(t *testing.T) {
	assert := assert.New(t)

	g := graph.NewGraph[int, string]()
	g.SetNode(1, "Node 1")
	g.SetNode(2, "Node 2")
	g.SetNode(3, "Node 3")
	g.AddEdge(3, 1)

	err := g.AddUndirectedEdge(1, 2)
	assert.Nil(err)

	node1Edges, err := g.GetEdges(1)
	assert.Nil(err)
	assert.Equal(1, len(node1Edges))

	node2Edges, err := g.GetEdges(2)
	assert.Nil(err)
	assert.Equal(1, len(node2Edges))

	err = g.AddUndirectedEdge(1, 2) // already defined
	assert.Nil(err)

	err = g.AddUndirectedEdge(3, 1) // already defined
	assert.Nil(err)
	err = g.AddUndirectedEdge(1, 3) // already defined
	assert.Nil(err)

	err = g.AddUndirectedEdge(1, 145) // non existing nodes
	assert.Error(err)

	err = g.AddUndirectedEdge(145, 1) // non existing nodes
	assert.Error(err)
}

func TestGetKey(t *testing.T) {
	assert := assert.New(t)

	type Person struct {
		name string
		age  int
	}

	g := graph.NewGraph[int, Person]()
	g.SetNode(1, Person{name: "Foo", age: 18})
	g.SetNode(2, Person{name: "Bar", age: 25})

	fooKey, ok := g.GetNodeKey(func(p Person) bool { return p.name == "Foo" })
	assert.Equal(1, fooKey)
	assert.True(ok)

	keys := g.GetNodeKeys(func(p Person) bool { return p.age >= 18 })
	assert.Equal(2, len(keys))

	_, ok = g.GetNodeKey(func(p Person) bool { return p.age == 45 })
	assert.False(ok)
}
