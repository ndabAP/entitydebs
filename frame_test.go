package entitydebs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ndabAP/entitydebs/tokenize"
)

func TestFrame_all(t *testing.T) {
	tests := []struct {
		name    string
		frame   frame
		offsets []int32
		tokens  [][]*tokenize.Token
	}{
		{
			name: "single sentence",
			frame: frame{
				sentences: []*tokenize.Sentence{
					{Text: &tokenize.TextSpan{BeginOffset: 0}},
				},
				tokens: []*tokenize.Token{
					{Text: &tokenize.TextSpan{BeginOffset: 0}},
					{Text: &tokenize.TextSpan{BeginOffset: 3}},
					{Text: &tokenize.TextSpan{BeginOffset: 6}},
				},
			},
			offsets: []int32{0},
			tokens: [][]*tokenize.Token{
				{
					{Text: &tokenize.TextSpan{BeginOffset: 0}},
					{Text: &tokenize.TextSpan{BeginOffset: 3}},
					{Text: &tokenize.TextSpan{BeginOffset: 6}},
				},
			},
		},
		{
			name: "multiple sentences",
			frame: frame{
				sentences: []*tokenize.Sentence{
					{Text: &tokenize.TextSpan{BeginOffset: 0}},
					{Text: &tokenize.TextSpan{BeginOffset: 10}},
				},
				tokens: []*tokenize.Token{
					{Text: &tokenize.TextSpan{BeginOffset: 0}},
					{Text: &tokenize.TextSpan{BeginOffset: 3}},
					{Text: &tokenize.TextSpan{BeginOffset: 6}},
					{Text: &tokenize.TextSpan{BeginOffset: 10}},
					{Text: &tokenize.TextSpan{BeginOffset: 12}},
				},
			},
			offsets: []int32{0, 3},
			tokens: [][]*tokenize.Token{
				{
					{Text: &tokenize.TextSpan{BeginOffset: 0}},
					{Text: &tokenize.TextSpan{BeginOffset: 3}},
					{Text: &tokenize.TextSpan{BeginOffset: 6}},
				},
				{
					{Text: &tokenize.TextSpan{BeginOffset: 10}},
					{Text: &tokenize.TextSpan{BeginOffset: 12}},
				},
			},
		},
		{
			name: "no sentences",
			frame: frame{
				sentences: []*tokenize.Sentence{},
				tokens:    []*tokenize.Token{},
			},
			offsets: []int32{},
			tokens:  [][]*tokenize.Token{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var i int32
			for offset, toks := range test.frame.all() {
				if diff := cmp.Diff(test.tokens[i], toks); diff != "" {
					t.Errorf("frame.all() mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(test.offsets[i], offset); diff != "" {
					t.Errorf("frame.all() mismatch (-want +got):\n%s", diff)
				}
				i++
			}
		})
	}
}
