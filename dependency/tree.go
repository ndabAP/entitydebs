package dependency

import (
	"github.com/ndabAP/entitydebs/tokenize"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

type (
	// Tree is a dependency tree.
	Tree struct {
		// A tree is a directed acyclic graph.
		graph *simple.DirectedGraph

		// root contains the root token.
		root *tokenize.Token
	}
)

var (
	_ graph.Node = (*tokenize.Token)(nil)
	_ dot.Node   = (*tokenize.Token)(nil)
)

// Root returns the root token.
func (tree Tree) Root() *tokenize.Token {
	return tree.root
}

// String returns a GraphViz DOT presentation of the tree.
func (tree Tree) String() string {
	b, _ := dot.Marshal(tree.graph, "", "", "	")
	return string(b)
}
