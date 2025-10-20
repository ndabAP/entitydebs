package entitydebs

import (
	"github.com/ndabAP/entitydebs/tokenize"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// Normalizer is a function that normalizes a token to reduce redundancy and
// improve data integrity. It is called with the token to normalize, its index
// within the frame, and all tokens of that frame.
//
// Note: Normalizers are not applied to entity tokens.
type Normalizer func(*tokenize.Token, int, []*tokenize.Token)

var (
	// NFKC applies Unicode Normalization Form KC.
	// See https://unicode.org/reports/tr15/#Norm_Forms. This is useful for
	// normalizing characters that look similar, e.g., "ﬀ" to "ff".
	NFKC Normalizer = func(token *tokenize.Token, _ int, _ []*tokenize.Token) {
		token.Text.Content = norm.NFKC.String(token.Text.Content)
	}

	// Lowercaser lowercases the token text. The language is set to
	// [language.Und] (undetermined). See ISO 639-2.
	Lowercaser Normalizer = func(token *tokenize.Token, _ int, _ []*tokenize.Token) {
		token.Text.Content = cases.Lower(language.Und).String(token.Text.Content)
	}

	// Lemma replaces the token text with its lemma. For example, "running" and
	// "ran" would both be normalized to "run".
	Lemma Normalizer = func(token *tokenize.Token, _ int, _ []*tokenize.Token) {
		token.Text.Content = token.Lemma
	}
)
