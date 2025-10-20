package dependency

import (
	"testing"

	"github.com/ndabAP/entitydebs/testhelper"
	"github.com/ndabAP/entitydebs/tokenize"
)

func TestTreeSearch(t *testing.T) {
	t.Parallel()

	sentence := testhelper.NewExampleTokens1(t, 0, 0)
	tree := Parse(0, sentence)

	tests := []struct {
		fn   func(token *tokenize.Token) bool
		want *tokenize.Token
		ok   bool
	}{
		{
			fn: func(token *tokenize.Token) bool {
				return token.Text.Content == sentence[0].Text.Content
			},
			want: sentence[0],
			ok:   true,
		},
		{
			fn: func(token *tokenize.Token) bool {
				return token.Text.Content == ""
			},
			want: &tokenize.Token{},
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			if got, ok := tree.Search(tt.fn); got != tt.want && ok != tt.ok {
				t.Errorf("Tree.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
