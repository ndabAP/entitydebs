package entitydebs

import (
	"slices"
	"strings"
)

type (
	// source wraps entities and texts, and returns a data [Frames].
	source struct {
		entity, texts []string
	}
)

// NewSource returns a new source, consisting of the entity, its aliases and
// texts. Duplicate entities and surrounding white spaces are removed.
//
// By convention, the first entity is the most well-known.
func NewSource(entity, texts []string) source {
	// De-duplicate entities
	dedup := make([]string, len(entity))
	for i, alias := range entity {
		if !slices.Contains(dedup, alias) {
			dedup[i] = alias
		}
	}

	// Trim leading and trailing white space.
	for i, alias := range dedup {
		dedup[i] = strings.TrimSpace(alias)
	}
	for i, text := range texts {
		texts[i] = strings.TrimSpace(text)
	}

	return source{
		entity: dedup,
		texts:  texts,
	}
}
