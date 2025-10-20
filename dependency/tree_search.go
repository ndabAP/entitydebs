package dependency

import "github.com/ndabAP/entitydebs/tokenize"

// Search iterates over all tokens and returns the found token and true if fn
// resolves to true. The token order is deterministic.
func (tree Tree) Search(fn func(token *tokenize.Token) bool) (*tokenize.Token, bool) {
	nodes := tree.graph.Nodes()
	for nodes.Next() {
		token := nodes.Node().(*tokenize.Token)
		if fn(token) {
			return token, true
		}
	}

	return nil, false
}
