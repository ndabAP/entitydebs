package dependency

import (
	"iter"

	"github.com/ndabAP/entitydebs/tokenize"
	"gonum.org/v1/gonum/graph/simple"
)

// Parse parses the given tokens and returns a head-to-dependent dependency
// tree. Parse will panic if the tokens contain cyclic dependencies.
//
// The offset can be used if the tokens' head indices cross sentence boundaries.
//
// Tokens and offset can't exceed 2^32.
func Parse(offset int32, tokens []*tokenize.Token) (tree Tree) {
	tree.graph = simple.NewDirectedGraph()
	for index, token := range tokensLimitedIter(tokens) {
		// We always start from the root token. A root token has no head tokens.
		if token.DependencyEdge.Label != tokenize.DependencyEdgeLabelRoot {
			continue
		}

		// Add the root token and parse the root tokens heads.
		tree.root = token
		tree.graph.AddNode(token)
		tree.parse(offset, index, tokens)

		break
	}

	return tree
}

// parse recursively parses tokens.
func (tree Tree) parse(offset, head int32, tokens []*tokenize.Token) {
	for index, token := range tokensLimitedIter(tokens) {
		// Ignore self.
		if index == head {
			continue
		}
		// We are looking for the tokens head.
		if token.DependencyEdge.HeadTokenIndex != head+offset {
			continue
		}

		// Add the token the dependent is headed by.
		func() {
			defer func() {
				if r := recover(); r != nil {
					panic("dependency: cyclic dependency")
				}
			}()
			tree.graph.AddNode(token)
		}()

		// Set the tokens dependency from head to dependent.
		edge := tree.graph.NewEdge(tokens[head], token)
		tree.graph.SetEdge(edge)

		// Parse the head tokens relations.
		tree.parse(offset, index, tokens)
	}
}

// tokensLimitedIter respects the 2^32 token limit.
func tokensLimitedIter(tokens []*tokenize.Token) iter.Seq2[int32, *tokenize.Token] {
	return func(yield func(int32, *tokenize.Token) bool) {
		for index, token := range tokens {
			if !yield(int32(index), token) {
				return
			}
		}
	}
}
