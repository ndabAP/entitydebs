package entitydebs

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ndabAP/entitydebs/dependency"
	"github.com/ndabAP/entitydebs/testhelper"
	"github.com/ndabAP/entitydebs/tokenize"
)

func TestFramesForest(t *testing.T) {
	t.Parallel()

	t.Run("one frame, entity in second sentence", func(t *testing.T) {
		t.Parallel()

		frames := Frames{}
		var (
			sentence1 = testhelper.NewExampleSentence1(t, 0)
			tokens1   = testhelper.NewExampleTokens1(t, 0, 0)

			charOffset = tokens1[len(tokens1)-1].Text.BeginOffset + 2 // Terminal plus white space.
			depOffset  = int32(len(tokens1)) - 1                      // 0-based index.
			sentence2  = testhelper.NewExampleSentence2(t, charOffset)
			tokens2    = testhelper.NewExampleTokens2(t, charOffset, depOffset)

			sentences = []*tokenize.Sentence{sentence1, sentence2}
			tokens    = slices.Concat(tokens1, tokens2)
			entities  = map[int][]*tokenize.Token{
				len(tokens1) + 5: {tokens2[5]}, // Housten
			}
		)
		frames.frames = []frame{
			{
				sentences: sentences,
				tokens:    tokens,
				entities:  entities,
			},
		}
		forest := frames.Forest()

		// Entites
		{
			var (
				got  = forest.entities[0]
				want = tokens2[5]
			)
			if got != want {
				t.Errorf("Frames.Forest().entities = %s, want %s", got, want)
			}
		}
		// Tokens
		{
			got := make([]*tokenize.Token, 0, len(tokens2))
			forest.Walk(func(token *tokenize.Token, _ dependency.Tree) bool {
				got = append(got, token)
				return true
			})
			// got order differs from tokens2.
			for _, tok := range got {
				if !slices.ContainsFunc(tokens2, func(t *tokenize.Token) bool {
					return t.ID() == tok.ID()
				}) {
					t.Fatalf("Frames.Forest().forest = %v, want %v", got, tokens2)
				}
			}
		}
	})
	t.Run("two frames, (multi token) entities in both", func(t *testing.T) {
		t.Parallel()

		var (
			entity1 = testhelper.NewToken(t,
				"",
				1,
				nil,
				testhelper.NewDepEdge(t, 0, tokenize.DependencyEdgeLabelUnknown),
				"",
			)
			entity2a = testhelper.NewToken(t,
				"",
				2,
				nil,
				testhelper.NewDepEdge(t, 0, tokenize.DependencyEdgeLabelUnknown),
				"",
			)
			entity2b = testhelper.NewToken(t,
				"",
				4,
				nil,
				testhelper.NewDepEdge(t, 2, tokenize.DependencyEdgeLabelUnknown),
				"",
			)
			toks1 = []*tokenize.Token{
				testhelper.NewToken(t,
					"",
					0,
					nil,
					testhelper.NewDepEdge(t, 0, tokenize.DependencyEdgeLabelRoot),
					"",
				),
				entity1,
				testhelper.NewToken(t,
					"",
					2,
					nil,
					testhelper.NewDepEdge(t, 1, tokenize.DependencyEdgeLabelUnknown),
					"",
				),
			}
			toks2 = []*tokenize.Token{
				testhelper.NewToken(t,
					"",
					3,
					nil,
					testhelper.NewDepEdge(t, 0, tokenize.DependencyEdgeLabelRoot),
					"",
				),
				testhelper.NewToken(t,
					"",
					4,
					nil,
					testhelper.NewDepEdge(t, 0, tokenize.DependencyEdgeLabelUnknown),
					"",
				),
			}
			toks3 = []*tokenize.Token{
				testhelper.NewToken(t,
					"",
					0,
					nil,
					testhelper.NewDepEdge(t, 1, tokenize.DependencyEdgeLabelUnknown),
					"",
				),
				testhelper.NewToken(t,
					"",
					1,
					nil,
					testhelper.NewDepEdge(t, 1, tokenize.DependencyEdgeLabelRoot),
					"",
				),
				entity2a,
				entity2b,
			}

			f = []frame{
				{
					sentences: []*tokenize.Sentence{
						testhelper.NewSentence(t,
							"",
							0,
							nil,
						),
						testhelper.NewSentence(t,
							"",
							3,
							nil,
						),
					},
					tokens: slices.Concat(toks1, toks2),
					entities: map[int][]*tokenize.Token{
						1: {entity1},
					},
				},
				{
					sentences: []*tokenize.Sentence{
						testhelper.NewSentence(t,
							"",
							0,
							nil,
						),
					},
					tokens: toks3,
					entities: map[int][]*tokenize.Token{
						2: {entity2a, entity2b},
					},
				},
			}
			frames = Frames{
				frames: f,
			}
		)

		forest := frames.Forest()
		roots := slices.Collect(forest.Roots())
		entities := slices.Collect(forest.Entities())
		heads := forest.Heads(nil)
		dependents := forest.Dependents(nil)

		if diff := cmp.Diff([]*tokenize.Token{entity1, entity2a, entity2b}, entities); diff != "" {
			t.Errorf("Frames.Forest().Entities() mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff([]*tokenize.Token{toks1[0], toks3[1]}, roots); diff != "" {
			t.Errorf("Frames.Forest().Roots() mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff([]*tokenize.Token{toks1[0], toks3[0]}, heads); diff != "" {
			t.Errorf("Frames.Forest().Heads() mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff([]*tokenize.Token{toks1[2]}, dependents); diff != "" {
			t.Errorf("Frames.Forest().Dependents() mismatch (-want +got):\n%s", diff)
		}
	})
}

func Test_depsHeads(t *testing.T) {
	t.Parallel()

	var (
		sentence1 = testhelper.NewExampleSentence1(t, 0)
		tokens1   = testhelper.NewExampleTokens1(t, 0, 0)
		_         = tokens1[0] // I
		tok1      = tokens1[1] // prefer
		tok2      = tokens1[2] // the
		tok3      = tokens1[3] // morning
		tok4      = tokens1[4] // flight
		tok5      = tokens1[5] // through
		tok6      = tokens1[6] // Denver
		tok7      = tokens1[7] // .

		charOffset = tok7.Text.BeginOffset + 2
		sentence2  = testhelper.NewExampleSentence2(t, charOffset)
		tokens2    = testhelper.NewExampleTokens2(t, charOffset, 0)
		_          = tokens2[0] // Book
		_          = tokens2[1] // me
		_          = tokens2[2] // the
		_          = tokens2[3] // flight
		tok12      = tokens2[4] // through
		tok13      = tokens2[5] // Houston
		_          = tokens2[6] // .
	)

	tests := []struct {
		name     string
		entities map[int][]*tokenize.Token
		fn       func(t *tokenize.Token) bool
		want     []*tokenize.Token
	}{
		{
			name: "single entity with head",
			entities: map[int][]*tokenize.Token{
				4: {tok4}, // flight
			},
			want: []*tokenize.Token{tok1}, // prefer
		},
		{
			name: "single entity with no heads",
			entities: map[int][]*tokenize.Token{
				1: {tok1}, // prefer
			},
			want: []*tokenize.Token{},
		},
		{
			name: "multi-token entity",
			entities: map[int][]*tokenize.Token{
				5: {tok5, tok6}, // through Denver
			},
			want: []*tokenize.Token{tok4}, // flight
		},
		{
			name: "entity in second frame",
			entities: map[int][]*tokenize.Token{
				13: {tok13}, // Housten
			},
			want: []*tokenize.Token{tok12}, // through
		},
		{
			name: "consecutive multi-token entity",
			entities: map[int][]*tokenize.Token{
				2: {tok2, tok3}, // the morning
			},
			want: []*tokenize.Token{tok4}, // flight
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			frames := Frames{
				frames: []frame{
					{
						sentences: []*tokenize.Sentence{sentence1},
						tokens:    tokens1,
						entities:  tt.entities, // All entities are pooled.
					},
					{
						sentences: []*tokenize.Sentence{sentence2},
						tokens:    tokens2,
					},
				},
			}

			got := frames.Forest().Heads(tt.fn)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("deps.Heads() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_depsDependents(t *testing.T) {
	t.Parallel()

	toks := testhelper.NewExampleTokens1(t, 0, 0)
	var (
		tok0 = toks[0] // I
		tok1 = toks[1] // prefer
		tok2 = toks[2] // the
		tok3 = toks[3] // morning
		tok4 = toks[4] // flight
		tok5 = toks[5] // through
		tok6 = toks[6] // Denver
		tok7 = toks[7] // .
	)
	sentence := &tokenize.Sentence{
		Text: &tokenize.TextSpan{
			Content: "I prefer the morning flight through Denver.",
		},
	}

	tests := []struct {
		name     string
		entities map[int][]*tokenize.Token
		want     []*tokenize.Token
	}{
		{
			name: "single entity with dependents",
			entities: map[int][]*tokenize.Token{
				4: {tok4}, // flight
			},
			want: []*tokenize.Token{tok2, tok3, tok5},
		},
		{
			name: "entity with no dependents",
			entities: map[int][]*tokenize.Token{
				7: {tok7}, // .
			},
			want: []*tokenize.Token{},
		},
		{
			name: "multi-token entity with nested dependents",
			entities: map[int][]*tokenize.Token{
				5: {tok5, tok6}, // through Denver
			},
			want: []*tokenize.Token{},
		},
		{
			name: "entity at root with dependents",
			entities: map[int][]*tokenize.Token{
				1: {tok1}, // prefer
			},
			want: []*tokenize.Token{tok0, tok4, tok7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			frames := Frames{
				frames: []frame{
					{
						sentences: []*tokenize.Sentence{sentence},
						tokens:    toks,
						entities:  tt.entities,
					},
				},
			}

			got := frames.Forest().Dependents(nil)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("deps.Dependents() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_depsRelationships(t *testing.T) {
	t.Parallel()

	toks := testhelper.NewExampleTokens1(t, 0, 0)
	var (
		tok0 = toks[0] // I
		tok1 = toks[1] // prefer
		tok2 = toks[2] // the
		tok3 = toks[3] // morning
		tok4 = toks[4] // flight
		tok5 = toks[5] // through
		tok6 = toks[6] // Denver
	)
	sentence := &tokenize.Sentence{
		Text: &tokenize.TextSpan{
			Content: "I prefer the morning flight through Denver.",
		},
	}

	tests := []struct {
		name     string
		entities map[int][]*tokenize.Token
		want     []*tokenize.Token
	}{
		{
			name: "single entity with heads and dependents",
			entities: map[int][]*tokenize.Token{
				4: {tok4},
			},
			want: []*tokenize.Token{tok1, tok2, tok3, tok5},
		},
		{
			name: "single entity with dependents",
			entities: map[int][]*tokenize.Token{
				6: {tok6},
			},
			want: []*tokenize.Token{tok5},
		},
		{
			name: "single entity with head",
			entities: map[int][]*tokenize.Token{
				0: {tok0}, // I
			},
			want: []*tokenize.Token{tok1}, // prefer
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			frames := Frames{
				frames: []frame{
					{
						sentences: []*tokenize.Sentence{sentence},
						tokens:    toks,
						entities:  tt.entities,
					},
				},
			}

			got := slices.Collect(frames.Forest().Relationships())
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("deps.Relations() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
