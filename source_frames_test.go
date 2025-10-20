package entitydebs

import (
	"bytes"
	"context"
	"slices"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
	"github.com/ndabAP/entitydebs/testhelper"
	"github.com/ndabAP/entitydebs/tokenize"
)

type mockTokenizer struct{}

func (tokenizer mockTokenizer) Tokenize(
	ctx context.Context,
	text string,
	feats tokenize.Features,
) (
	tokenize.Analysis,
	error,
) {
	var (
		tokens    []*tokenize.Token
		sentences []*tokenize.Sentence
		sentiment *tokenize.Sentiment
	)

	if feats&tokenize.FeatureSyntax != 0 {
		t, s := tokenizer.syntax(text)
		tokens = t
		sentences = s
	}
	if feats&tokenize.FeatureSentiment != 0 {
		sentiment = &tokenize.Sentiment{
			Score:     -1.0,
			Magnitude: 1.0,
		}
	}

	analysis := tokenize.Analysis{
		Tokens:    tokens,
		Sentences: sentences,
		Sentiment: sentiment,
	}
	return analysis, nil
}

func (tokenizer mockTokenizer) syntax(text string) ([]*tokenize.Token, []*tokenize.Sentence) {
	var (
		tokens    = make([]*tokenize.Token, 0)
		sentences = make([]*tokenize.Sentence, 0)

		token    = &strings.Builder{}
		sentence = &strings.Builder{}

		runes = utf8.RuneCountInString(text)
	)
	for i := 1; i <= runes; i++ {
		r, _ := utf8.DecodeRuneInString(text[i-1:])

		// White space
		if unicode.Is(unicode.Space, r) {
			// Terminal token
			var tok string

			// Special case: Preceding token is sentence terminal.
			if token.Len() > 0 {
				tok, tokens = tokenizer.terminate(token, tokens)
			}

			// Append token to sentence.
			sentence.WriteString(tok)

			// Special case: Sentences shan't beginn with white space.
			if sentence.Len() > 0 {
				sentence.WriteRune(r)
			}

			continue
		}

		// ., !, ?
		if unicode.Is(unicode.Sentence_Terminal, r) {
			var tok string
			tok, tokens = tokenizer.terminate(token, tokens)

			// Append terminal token to tokens
			tokens = append(tokens, &tokenize.Token{Text: &tokenize.TextSpan{Content: string(r)}})

			// Terminate sentence
			sentence.WriteString(tok)
			sentence.WriteRune(r)
			sentences = append(sentences, &tokenize.Sentence{Text: &tokenize.TextSpan{Content: sentence.String()}})
			sentence.Reset()

			continue
		}

		// A-Za-z, 0-9
		if unicode.In(r, unicode.Letter, unicode.Number) {
			token.WriteRune(r)
		} else {
			// Not supported, using replacement character.
			token.WriteRune(unicode.ReplacementChar)
		}

		// Special case: Final token is not sentence terminal.
		if i == runes {
			var tok string
			tok, tokens = tokenizer.terminate(token, tokens)

			// Terminate sentence
			sentence.WriteString(tok)
			sentences = append(sentences, &tokenize.Sentence{Text: &tokenize.TextSpan{Content: sentence.String()}})
			break
		}
	}

	return tokens, sentences
}

func (tokenizer mockTokenizer) terminate(
	token *strings.Builder,
	tokens []*tokenize.Token,
) (
	string,
	[]*tokenize.Token,
) {
	tok := token.String()
	tokens = append(tokens, &tokenize.Token{Text: &tokenize.TextSpan{Content: tok}})
	token.Reset()
	return tok, tokens
}

func Test_sourceFrames(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	type (
		want struct {
			entities map[string][]tokenize.Token
			frames   []frame
		}
		test struct {
			source     source
			feats      tokenize.Features
			normalizer []Normalizer
			want       want
			err        error
		}
	)

	tests := make([]test, 0)
	{
		var (
			tokens = []*tokenize.Token{
				testhelper.NewToken(t, "Everything", 0, nil, nil, ""),
				testhelper.NewToken(t, "ripped", 0, nil, nil, ""),
				testhelper.NewToken(t, "apart", 0, nil, nil, ""),
				testhelper.NewToken(t, "in", 0, nil, nil, ""),
				testhelper.NewToken(t, "a", 0, nil, nil, ""),
				testhelper.NewToken(t, "New", 0, nil, nil, ""),
				testhelper.NewToken(t, "York", 0, nil, nil, ""),
				testhelper.NewToken(t, "minute", 0, nil, nil, ""),
				testhelper.NewToken(t, ".", 0, nil, nil, ""),
			}
			texts = []string{
				"Everything ripped apart in a New York minute.",
			}
			entity = []string{
				"New York",
			}
		)
		tests = append(tests, test{
			source: source{
				texts:  texts,
				entity: entity,
			},
			feats: tokenize.FeatureSyntax,
			want: want{
				entities: map[string][]tokenize.Token{
					entity[0]: {
						*tokens[5],
						*tokens[6],
					},
				},
				frames: []frame{
					{
						sentences: []*tokenize.Sentence{
							testhelper.NewSentence(t,
								texts[0],
								0,
								nil,
							),
						},
						tokens: tokens,
						entities: map[int][]*tokenize.Token{
							5: {
								tokens[5],
								tokens[6],
							},
						},
					},
				},
			},
		})
	}
	{
		var (
			tokens = []*tokenize.Token{
				testhelper.NewToken(t, "Pissing", 0, nil, nil, ""),
				testhelper.NewToken(t, "Punchinello", 0, nil, nil, ""),
				testhelper.NewToken(t, "off", 0, nil, nil, ""),
				testhelper.NewToken(t, "was", 0, nil, nil, ""),
				testhelper.NewToken(t, "a", 0, nil, nil, ""),
				testhelper.NewToken(t, "dangerous", 0, nil, nil, ""),
				testhelper.NewToken(t, "game", 0, nil, nil, ""),
			}
			texts = []string{
				"Pissing Punchinello off was a dangerous game",
			}
			entity = []string{
				"Punchinello",
			}
		)
		tests = append(tests, test{
			source: source{
				texts:  texts,
				entity: entity,
			},
			feats: tokenize.FeatureSyntax,
			want: want{
				entities: map[string][]tokenize.Token{
					entity[0]: {*tokens[1]},
				},
				frames: []frame{
					{
						sentences: []*tokenize.Sentence{
							testhelper.NewSentence(t,
								texts[0],
								0,
								nil,
							),
						},
						tokens: tokens,
						entities: map[int][]*tokenize.Token{
							1: {tokens[1]},
						},
					},
				},
			},
		})
	}
	{
		var (
			tokens = []*tokenize.Token{
				testhelper.NewToken(t, "Punchinello", 0, nil, nil, ""),
				testhelper.NewToken(t, "was", 0, nil, nil, ""),
				testhelper.NewToken(t, "a", 0, nil, nil, ""),
				testhelper.NewToken(t, "pushover", 0, nil, nil, ""),
				testhelper.NewToken(t, ".", 0, nil, nil, ""),
				testhelper.NewToken(t, "The", 0, nil, nil, ""),
				testhelper.NewToken(t, "moment", 0, nil, nil, ""),
				testhelper.NewToken(t, "I", 0, nil, nil, ""),
				testhelper.NewToken(t, "stepped", 0, nil, nil, ""),
				testhelper.NewToken(t, "into", 0, nil, nil, ""),
				testhelper.NewToken(t, "the", 0, nil, nil, ""),
				testhelper.NewToken(t, "room", 0, nil, nil, ""),
				testhelper.NewToken(t, "he", 0, nil, nil, ""),
				testhelper.NewToken(t, "folded", 0, nil, nil, ""),
				testhelper.NewToken(t, "like", 0, nil, nil, ""),
				testhelper.NewToken(t, "a", 0, nil, nil, ""),
				testhelper.NewToken(t, "deuce", 0, nil, nil, ""),
				testhelper.NewToken(t, "before", 0, nil, nil, ""),
				testhelper.NewToken(t, "a", 0, nil, nil, ""),
				testhelper.NewToken(t, "royal", 0, nil, nil, ""),
				testhelper.NewToken(t, "flush", 0, nil, nil, ""),
				testhelper.NewToken(t, ".", 0, nil, nil, ""),
			}
			texts = []string{
				"Punchinello was a pushover. The moment I stepped into the room he folded like a deuce before a royal flush.",
			}
			entity = []string{
				"Punchinello",
			}
		)
		tests = append(tests, test{
			source: source{
				texts:  texts,
				entity: entity,
			},
			feats: tokenize.FeatureSyntax,
			want: want{
				entities: map[string][]tokenize.Token{
					entity[0]: {*tokens[0]},
				},
				frames: []frame{
					{
						sentences: []*tokenize.Sentence{
							{Text: &tokenize.TextSpan{Content: "Punchinello was a pushover."}},
							{Text: &tokenize.TextSpan{Content: "The moment I stepped into the room he folded like a deuce before a royal flush."}},
						},
						tokens: tokens,
						entities: map[int][]*tokenize.Token{
							0: {tokens[0]},
						},
					},
				},
			},
		})
	}
	{
		var (
			tokens = []*tokenize.Token{
				testhelper.NewToken(t, "V", 0, nil, nil, ""),
				testhelper.NewToken(t, "for", 0, nil, nil, ""),
				testhelper.NewToken(t, "Valhalla", 0, nil, nil, ""),
				testhelper.NewToken(t, ".", 0, nil, nil, ""),
			}
			texts = []string{
				"V for Valhalla.",
			}
			entity = []string{
				"V",
				"Valhalla",
			}
		)
		tests = append(tests, test{
			source: source{
				texts:  texts,
				entity: entity,
			},
			feats: tokenize.FeatureSyntax,
			want: want{
				entities: map[string][]tokenize.Token{
					entity[0]: {*tokens[0]},
					entity[1]: {*tokens[2]},
				},
				frames: []frame{
					{
						sentences: []*tokenize.Sentence{
							testhelper.NewSentence(t,
								texts[0],
								0,
								nil,
							),
						},
						tokens: tokens,
						entities: map[int][]*tokenize.Token{
							0: {tokens[0]},
							2: {tokens[2]},
						},
					},
				},
			},
		})
	}
	{
		var (
			char    = string(unicode.MaxRune)
			invalid = string(unicode.ReplacementChar)

			tokens = []*tokenize.Token{
				testhelper.NewToken(t, invalid, 0, nil, nil, ""),
			}
			texts     = []string{char}
			entity    = []string{}
			sentiment = &tokenize.Sentiment{
				Score:     -1.0,
				Magnitude: 1.0,
			}
		)
		tests = append(tests, test{
			source: source{
				texts:  texts,
				entity: entity,
			},
			feats: tokenize.FeatureSyntax | tokenize.FeatureSentiment,
			want: want{
				entities: map[string][]tokenize.Token{},
				frames: []frame{
					{
						sentences: []*tokenize.Sentence{
							testhelper.NewSentence(t,
								invalid,
								0,
								nil,
							),
						},
						tokens: []*tokenize.Token{
							tokens[0],
						},
						sentiment: sentiment,
						entities:  map[int][]*tokenize.Token{},
					},
				},
			},
		})
	}
	{
		var (
			tokens = []*tokenize.Token{
				testhelper.NewToken(t, "Everything", 0, nil, nil, ""),
				testhelper.NewToken(t, "ripped", 0, nil, nil, ""),
				testhelper.NewToken(t, "apart", 0, nil, nil, ""),
				testhelper.NewToken(t, "in", 0, nil, nil, ""),
				testhelper.NewToken(t, "a", 0, nil, nil, ""),
				testhelper.NewToken(t, "New", 0, nil, nil, ""),
				testhelper.NewToken(t, "York", 0, nil, nil, ""),
				testhelper.NewToken(t, "minute", 0, nil, nil, ""),
				testhelper.NewToken(t, ".", 0, nil, nil, ""),
			}
			texts = []string{
				"Everything ripped apart in a New York minute.",
			}
			entity = []string{
				"New York",
			}
		)
		tests = append(tests, test{
			source: source{
				texts:  texts,
				entity: entity,
			},
			feats: tokenize.FeatureSyntax,
			normalizer: []Normalizer{
				Lowercaser,
			},
			want: want{
				entities: map[string][]tokenize.Token{
					entity[0]: {
						*tokens[5],
						*tokens[6],
					},
				},
				frames: []frame{
					{
						sentences: []*tokenize.Sentence{
							{Text: &tokenize.TextSpan{Content: "Everything ripped apart in a New York minute."}},
						},
						tokens: tokens,
						entities: map[int][]*tokenize.Token{
							5: {
								tokens[5],
								tokens[6],
							},
						},
					},
				},
			},
		})
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			tokenizer := mockTokenizer{}
			frames, err := test.source.Frames(
				ctx,
				tokenizer,
				test.feats,
				test.normalizer...,
			)
			if err != test.err {
				t.Fatalf("source.Frames(%s, %s, %d) = _, %s, want %s",
					ctx,
					"mockTokenizer{}",
					test.feats,
					err,
					test.err,
				)
			}

			normalize(t, test.want.frames, nil, test.normalizer...)
			got, want := testhelper.MarshalJSON(t, frames), testhelper.MarshalJSON(t, Frames{
				frames:   test.want.frames,
				entities: test.want.entities,
			})
			if !bytes.Equal(got, want) {
				t.Errorf("source.Frames(%s, %s, %d) = %s, _, want %s ",
					ctx,
					"mockTokenizer{}",
					test.feats,
					got,
					want,
				)
				t.Error(cmp.Diff(got, want))
			}
		})
	}
}

func normalize(t *testing.T, frames []frame, tokens []*tokenize.Token, normalizer ...Normalizer) {
	t.Helper()

	for _, frame := range frames {
		s := make([]int, 0)
		for i, alias := range frame.entities {
			s = append(s, i)
			for j := range alias[1:] {
				s = append(s, j+i+1)
			}
		}

		for i, token := range frame.tokens {
			if slices.Contains(s, i) {
				continue
			}

			for i, f := range normalizer {
				f(token, i, tokens)
			}
		}
	}
}
