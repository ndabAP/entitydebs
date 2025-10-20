package entitydebs

import (
	"context"
	"iter"
	"maps"
	"slices"

	"github.com/ndabAP/entitydebs/tokenize"
)

// Frames tokenizes all texts and entities of source into [Frames]
// according to tokenizer. An [Frames] is a collection of data frames and
// entites within data frames.
//
// [Normalizer] can be used to to reduce redundancy and improve data integrity.
// Normalizers are not applied to entity tokens.
func (source source) Frames(
	ctx context.Context,
	tokenizer tokenize.Tokenizer,
	feats tokenize.Features,
	normalizer ...Normalizer,
) (
	frames Frames,
	err error,
) {
	// Tokenize entities.
	entities := make(map[string][]tokenize.Token, len(source.entity))
	for _, entity := range source.entity {
		select {
		case <-ctx.Done():
			return frames, ctx.Err()
		default:
		}

		analysis, err := tokenizer.Tokenize(ctx, entity, tokenize.FeatureSyntax)
		if err != nil {
			return frames, err
		}
		for _, token := range analysis.Tokens {
			entities[entity] = append(entities[entity], *token.Clone())
		}
	}
	frames.entities = entities

	// Tokenize texts into data frames.
	frames.frames = make([]frame, 0, len(source.texts))
	for _, text := range source.texts {
		select {
		case <-ctx.Done():
			return frames, ctx.Err()
		default:
		}

		f, err := source.frame(ctx, tokenizer, text, entities, feats, normalizer...)
		if err != nil {
			return frames, err
		}
		frames.frames = append(frames.frames, f)

	}

	return
}

// frame computes a single data frame.
func (source source) frame(
	ctx context.Context,
	tokenizer tokenize.Tokenizer,
	text string,
	entities map[string][]tokenize.Token,
	feats tokenize.Features,
	normalizer ...Normalizer,
) (
	frame frame,
	err error,
) {
	analysis, err := tokenizer.Tokenize(ctx, text, feats)
	if err != nil {
		return
	}

	frame.tokens = make([]*tokenize.Token, len(analysis.Tokens))
	frame.sentences = make([]*tokenize.Sentence, len(analysis.Sentences))
	frame.sentences = analysis.Sentences
	frame.entities = make(map[int][]*tokenize.Token, 0)
	frame.sentiment = analysis.Sentiment

	i := 0
	for i != len(analysis.Tokens) {
		// Peek for entity tokens.
		_, tokens, j := peek(analysis.Tokens[i:], entities)
		switch j {
		// No entity tokens found.
		case -1:
			// Append only to tokens.
			token := analysis.Tokens[i]
			frame.tokens[i] = token

			// Apply normalizer, if any.
			for _, f := range normalizer {
				f(token, i, analysis.Tokens)
			}

			i++

		// Found one or more entity tokens.
		default:
			for k, t := range tokens {
				// Append entity tokens to both, tokens and entities.
				frame.tokens[i+k] = t
				frame.entities[i] = append(frame.entities[i], t)
			}

			// Skip about entity positions.
			i += len(tokens)
		}
	}

	return
}

// peek checks if subsequent tokens are entity tokens. It returns the entity as
// string, as tokens and the final index. If no entity was found, peek returns
// -1.
func peek(tokens []*tokenize.Token, entities map[string][]tokenize.Token) (string, []*tokenize.Token, int) {
	// i contains the final entity index.
	i := 0

	// Entity alias iterator
	next, stop := iter.Pull2(maps.All(entities))
	defer stop()
	for {
		found := true
		entity, toks, ok := next()
		if !ok {
			break
		}

		// Entity buffer
		buf := make([]*tokenize.Token, 0, len(toks))

		// Entity iterator
		n, s := iter.Pull2(slices.All(toks))
		// Text iterator
		m, t := iter.Pull2(slices.All(tokens))
		for {
			// If no entity is left, cancel.
			j, v, ok := n()
			if !ok {
				s()
				t()
				break
			}
			// If no text is left, cancel.
			_, w, ok := m()
			if !ok {
				s()
				t()
				break
			}

			if w.Text.Content != v.Text.Content {
				i = 0
				found = false
				s()
				t()
				// Continue with next entity.
				break
			}

			buf = append(buf, w)
			i = j
		}

		if found {
			return entity, buf, i
		}
	}

	return "", nil, -1
}
