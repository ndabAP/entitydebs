package entitydebs

import (
	"iter"

	"github.com/ndabAP/entitydebs/tokenize"
)

// frame represents a text, consisting of sentences, tokens, sentiment and
// entities.
type frame struct {
	sentences []*tokenize.Sentence
	tokens    []*tokenize.Token
	sentiment *tokenize.Sentiment

	// entities is a map of the entities starting offset and consecutive tokens.
	// The final offset can be obtained by the adding the offset to the length
	// of the tokens. The offset is zero-based.
	entities map[int][]*tokenize.Token
}

// all yields all tokens of all sentences of a frame and the token offset to the
// first token of a sentence.
func (f frame) all() iter.Seq2[int32, []*tokenize.Token] {
	var (
		sentences = f.sentences
		toks      = f.tokens

		// j is the total token offset.
		j int32
		// offset of the current sentence.
		offset int32
	)
	return func(yield func(int32, []*tokenize.Token) bool) {
		for i := range sentences {
			// next is the next sentence offset, if any.
			var next int32 = -1
			if i < len(sentences)-1 {
				next = sentences[i+1].Text.BeginOffset
			}

			tokens := make([]*tokenize.Token, 0, len(toks[j:]))
			for _, token := range toks[j:] {
				if next == token.Text.BeginOffset {
					// Next sentence
					break
				}

				j++
				tokens = append(tokens, token)
			}
			if !yield(offset, tokens) {
				break
			}

			offset = j
		}
	}
}
