package entitydebs

import (
	"iter"

	"github.com/ndabAP/entitydebs/tokenize"
)

// Frames contains all analyzed frames and entities of a source. A forest of
// dependency trees can be constructed from the frames.
type Frames struct {
	// frames contains data frames for each analyzed text.
	frames []frame

	// entities contains the tokenized entity.
	entities map[string][]tokenize.Token

	// deps contains dependency trees of all frames.
	deps deps
}

// All returns all tokens of each sentence of all frames and the token offset to
// the first token of a sentence.
//
// An offset of zero indicates a frame advancement.
func (f Frames) All() iter.Seq2[int32, []*tokenize.Token] {
	return func(yield func(int32, []*tokenize.Token) bool) {
		for _, frame := range f.frames {
			for i, sentence := range frame.all() {
				if !yield(i, sentence) {
					return
				}
			}
		}
	}
}
