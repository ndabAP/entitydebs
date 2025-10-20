package entitydebs

import (
	"iter"
	"math"
	"slices"

	"github.com/ndabAP/entitydebs/dependency"
	"github.com/ndabAP/entitydebs/tokenize"
)

type (
	deps struct {
		// forest contains all dependency forest of all frames.
		forest []dependency.Tree

		// entities contains all entity tokens of all frames.
		entities []*tokenize.Token
	}
)

// Forest returns a forest of dependency trees that contain entity tokens.
// The first call to Forest constructs the forest, subsequent calls return the
// cached forest.
func (f *Frames) Forest() deps {
	if f.deps.forest != nil {
		return f.deps
	}

	deps := deps{
		forest:   make([]dependency.Tree, 0),
		entities: make([]*tokenize.Token, 0),
	}

	// Accumulate all entities of all frames.
	for _, frame := range f.frames {
		for _, tokens := range frame.entities {
			deps.entities = append(deps.entities, tokens...)
		}
	}

	// Accumulate all dependency trees of all frames.
	var offset int32
	for index, tokens := range f.All() {
		if index == 0 {
			// Next frame, reset offset.
			offset = 0
		}

		// We are looking for trees that contain entity tokens.
		if !slices.ContainsFunc(tokens, func(token *tokenize.Token) bool {
			return slices.Contains(deps.entities, token)
		}) {
			// Don't corrupt offset.
			n := len(tokens)
			if n > math.MaxInt32 {
				panic("entitydebs: too many tokens")
			}
			offset += int32(n)

			continue
		}
		offset += index

		tree := dependency.Parse(offset, tokens)
		deps.forest = append(deps.forest, tree)
	}
	f.deps = deps

	return f.deps
}

// Trees returns all dependency trees of all frames.
func (f Frames) Trees() iter.Seq[dependency.Tree] {
	return func(yield func(dependency.Tree) bool) {
		for _, tree := range f.Forest().forest {
			if !yield(tree) {
				return
			}
		}
	}
}

// Roots returns all root tokens for every tree.
func (deps deps) Roots() iter.Seq[*tokenize.Token] {
	return func(yield func(*tokenize.Token) bool) {
		for _, tree := range deps.forest {
			if !yield(tree.Root()) {
				return
			}
		}
	}
}

// Entities returns all entities of all trees.
func (deps deps) Entities() iter.Seq[*tokenize.Token] {
	return func(yield func(*tokenize.Token) bool) {
		for _, entity := range deps.entities {
			if !yield(entity) {
				return
			}
		}
	}
}

// Heads returns all head tokens from entity tokens. Every token is called
// with fn. If fn returns false, Heads is cancelled.
func (deps deps) Heads(fn func(t *tokenize.Token) bool) []*tokenize.Token {
	var (
		heads = make([]*tokenize.Token, 0)
		seen  = make(map[*tokenize.Token]struct{})
	)
	deps.walker(func(token *tokenize.Token, tree dependency.Tree) bool {
		if !slices.Contains(deps.entities, token) {
			return true
		}

		// Skip connected entity tokens.
		var (
			dependent = token
			head      *tokenize.Token
		)
		for {
			head = tree.Head(dependent)
			if head == nil {
				// No head found.
				return true
			}
			if _, ok := seen[head]; ok {
				// Head already seen, next token.
				return true
			}
			if !slices.Contains(deps.entities, head) {
				// Final head found.
				break
			}

			// Retry with head of head.
			dependent = head
		}

		// Head found.
		heads = append(heads, head)
		if fn != nil && !fn(head) {
			return false
		}
		seen[head] = struct{}{}

		return true
	})

	return heads
}

// Dependents returns all dependent tokens from entity tokens.
func (deps deps) Dependents(fn func(t *tokenize.Token) bool) []*tokenize.Token {
	dependents := make([]*tokenize.Token, 0)
	deps.walker(func(token *tokenize.Token, tree dependency.Tree) bool {
		if !slices.Contains(deps.entities, token) {
			return true
		}

		head := token
		for !tree.Dependents(head, func(dependent *tokenize.Token) bool {
			// Follow multi-token entities.
			if slices.Contains(deps.entities, dependent) {
				// Retry with dependent of dependent.
				head = dependent
				return false
			}

			// Dependent found.
			dependents = append(dependents, dependent)
			if fn != nil && !fn(dependent) {
				return false
			}

			return true
		}) {
			// No more dependents.
			if head == token {
				break
			}
		}

		return true
	})

	return dependents
}

// Relationships returns all heads and dependents tokens from entity tokens.
func (deps deps) Relationships() iter.Seq[*tokenize.Token] {
	return func(yield func(*tokenize.Token) bool) {
		for _, head := range deps.Heads(nil) {
			if !yield(head) {
				return
			}
		}
		for _, dependent := range deps.Dependents(nil) {
			if !yield(dependent) {
				return
			}
		}
	}
}

// Dependencies iterates over all dependencies, with fn receiving the head,
// dependent, dependency edge label, and tree.
func (deps deps) Dependencies(
	fn func(
		head,
		dependent *tokenize.Token,
		dep tokenize.DependencyEdgeLabel,
		tree dependency.Tree,
	) bool,
) {
	for _, tree := range deps.forest {
		tree.Dependencies(func(head, dependent *tokenize.Token) bool {
			dep := dependent.DependencyEdge.Label
			return fn(head, dependent, dep, tree)
		})
	}
}

// Walk walks every token of every tree in a forest in pre-order and executes
// fn.
func (deps deps) Walk(fn func(token *tokenize.Token, tree dependency.Tree) bool) {
	deps.walker(fn)
}

// walker walks every token of every tree in a forest and executes fn.
func (deps deps) walker(fn func(token *tokenize.Token, tree dependency.Tree) bool) {
	for _, tree := range deps.forest {
		tree.Walk(func(token *tokenize.Token) bool { return fn(token, tree) })
	}
}
