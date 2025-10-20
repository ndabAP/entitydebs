package dependency

import (
	"cmp"
	"slices"

	"github.com/ndabAP/entitydebs/tokenize"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/traverse"
)

var _ traverse.Graph = (*Tree)(nil)

// From implements the graph traverser.
func (tree Tree) From(id int64) graph.Nodes {
	tokens := make([]graph.Node, 0)
	iter := tree.graph.From(id)
	if iter == graph.Empty {
		return graph.Empty
	}
	for iter.Next() {
		tokens = append(tokens, iter.Node())
	}

	// Sort the tokens in descending order.
	slices.SortFunc(tokens, func(a, b graph.Node) int {
		return cmp.Compare(b.ID(), a.ID())
	})

	return iterator.NewOrderedNodes(tokens)
}

func (tree Tree) To(id int64) graph.Nodes {
	tokens := make([]graph.Node, 0)
	iter := tree.graph.To(id)
	if iter == graph.Empty {
		return graph.Empty
	}
	for iter.Next() {
		tokens = append(tokens, iter.Node())
	}

	// Sort the tokens in descending order.
	slices.SortFunc(tokens, func(a, b graph.Node) int {
		return cmp.Compare(b.ID(), a.ID())
	})

	return iterator.NewOrderedNodes(tokens)
}

// Edge implements the graph traverser.
func (tree Tree) Edge(uid, vid int64) graph.Edge {
	return tree.graph.Edge(uid, vid)
}

// Walk walks the tree in pre-order and calls fn on every visited node. It can
// be cancelled with returning false.
func (tree Tree) Walk(fn func(t *tokenize.Token) bool) {
	traverser := traverse.DepthFirst{}
	var (
		from  = tree.Root()
		until = func(node graph.Node) bool {
			return !fn(node.(*tokenize.Token))
		}
	)
	traverser.Walk(
		tree,
		from,
		until,
	)
}

// Dependencies iterates over every relationship and calls fn with the head and
// dependent token. If fn returns false, Dependencies returns.
func (tree Tree) Dependencies(fn func(head, dependent *tokenize.Token) bool) {
	tokens := make([]graph.Node, 0)
	iter := tree.graph.Nodes()
	if iter == graph.Empty {
		return
	}
	for iter.Next() {
		tokens = append(tokens, iter.Node())
	}
	// Sort the tokens in ascending order.
	slices.SortFunc(tokens, func(a, b graph.Node) int {
		return cmp.Compare(a.ID(), b.ID())
	})

	for _, node := range tokens {
		head := node.(*tokenize.Token)
		nodes := tree.From(head.ID())
		for nodes.Next() {
			dependent := nodes.Node().(*tokenize.Token)
			if !fn(head, dependent) {
				return
			}
		}
	}
}

// Head returns the token dependent is headed from. If dependent is the root
// token, Head returns nil.
func (tree Tree) Head(dependent *tokenize.Token) *tokenize.Token {
	iter := tree.graph.To(dependent.ID())
	if iter == graph.Empty {
		return nil
	}

	var head *tokenize.Token
	for iter.Next() {
		head = iter.Node().(*tokenize.Token)
	}
	return head
}

// Dependents iterates over every dependent of head and calls fn with the token.
// If head has no dependents, it returns false, otherwise true.
func (tree Tree) Dependents(head *tokenize.Token, fn func(*tokenize.Token) bool) bool {
	tokens := make([]graph.Node, 0)
	iter := tree.graph.From(head.ID())
	if iter == graph.Empty {
		return false
	}
	for iter.Next() {
		tokens = append(tokens, iter.Node())
	}
	// Sort the tokens in ascending order.
	slices.SortFunc(tokens, func(a, b graph.Node) int {
		return cmp.Compare(a.ID(), b.ID())
	})

	iter = iterator.NewOrderedNodes(tokens)
	walker(iter, fn)
	return true
}

// walker is a helper function which takes a token iterator and returns if fn
// resolves to false.
func walker(nodes graph.Nodes, fn func(token *tokenize.Token) bool) {
	for nodes.Next() {
		token := nodes.Node().(*tokenize.Token)
		if !fn(token) {
			return
		}
	}
}
