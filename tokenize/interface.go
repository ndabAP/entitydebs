package tokenize

import (
	"context"
)

// Tokenizer is an interface for tokenizing text with the specified features.
type Tokenizer interface {
	Tokenize(ctx context.Context, text string, feats Features) (Analysis, error)
}
