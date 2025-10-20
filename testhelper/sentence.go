package testhelper

import (
	"testing"

	"github.com/ndabAP/entitydebs/tokenize"
)

func NewExampleSentence1(t *testing.T, offset int32) *tokenize.Sentence {
	t.Helper()

	return &tokenize.Sentence{
		Text: &tokenize.TextSpan{
			Content:     "I prefer the morning flight through Denver.",
			BeginOffset: offset,
		},
		Sentiment: &tokenize.Sentiment{
			Magnitude: 0.10400000214576721,
			Score:     0.04100000113248825,
		},
	}
}

func NewExampleSentence2(t *testing.T, offset int32) *tokenize.Sentence {
	t.Helper()

	return &tokenize.Sentence{
		Text: &tokenize.TextSpan{
			Content:     "Book me the flight through Houston.",
			BeginOffset: offset,
		},
		Sentiment: &tokenize.Sentiment{
			Magnitude: 0.039000,
			Score:     0.006000,
		},
	}
}
