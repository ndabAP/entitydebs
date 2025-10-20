package entitydebs

import (
	"testing"

	"github.com/ndabAP/entitydebs/tokenize"
)

func TestNFKC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		token *tokenize.Token
		want  string
	}{
		{
			name: "ligature ff",
			token: &tokenize.Token{
				Text: &tokenize.TextSpan{Content: "ﬀ"},
			},
			want: "ff",
		},
		{
			name: "no change",
			token: &tokenize.Token{
				Text: &tokenize.TextSpan{Content: "hello"},
			},
			want: "hello",
		},
		{
			name: "full-width to half-width",
			token: &tokenize.Token{
				Text: &tokenize.TextSpan{Content: "ＨＥＬＬＯ"},
			},
			want: "HELLO",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			NFKC(tt.token, 0, nil)
			if tt.token.Text.Content != tt.want {
				t.Errorf("NFKC() got = %q, want %q", tt.token.Text.Content, tt.want)
			}
		})
	}
}

func TestLowercaser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "English mixed case",
			input: "HeLlO wOrLd",
			want:  "hello world",
		},
		{
			name:  "German sharp S",
			input: "Straße",
			want:  "straße",
		},
		{
			name:  "Turkish dotted I",
			input: "İSTANBUL",
			want:  "i̇stanbul",
		},
		{
			name:  "Greek mixed case",
			input: "Αθήνα",
			want:  "αθήνα",
		},
		{
			name:  "Russian mixed case",
			input: "Москва",
			want:  "москва",
		},
		{
			name:  "Japanese Katakana full-width",
			input: "カタカナ",
			want:  "カタカナ",
		},
		{
			name:  "Undetermined language",
			input: "HeLlO wOrLd",
			want:  "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			token := &tokenize.Token{
				Text: &tokenize.TextSpan{Content: tt.input},
			}
			Lowercaser(token, 0, nil)
			if token.Text.Content != tt.want {
				t.Errorf(
					"Lowercaser() got = %q, want %q",
					tt.input,
					tt.want,
				)
			}
		})
	}
}

func TestLemma(t *testing.T) {
	tests := []struct {
		name  string
		token *tokenize.Token
		want  string
	}{
		{
			name: "running to run",
			token: &tokenize.Token{
				Text:  &tokenize.TextSpan{Content: "running"},
				Lemma: "run",
			},
			want: "run",
		},
		{
			name: "ran to run",
			token: &tokenize.Token{
				Text:  &tokenize.TextSpan{Content: "ran"},
				Lemma: "run",
			},
			want: "run",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			Lemma(tt.token, 0, nil)
			if tt.token.Text.Content != tt.want {
				t.Errorf("Lemma() got = %q, want %q", tt.token.Text.Content, tt.want)
			}
		})
	}
}
