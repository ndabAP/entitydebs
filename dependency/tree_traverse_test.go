package dependency

import (
	"slices"
	"testing"

	"github.com/ndabAP/entitydebs/testhelper"
	"github.com/ndabAP/entitydebs/tokenize"
)

func TestTreeWalk(t *testing.T) {
	t.Parallel()

	sentence := testhelper.NewExampleTokens1(t, 0, 0)

	want := []string{
		"prefer",
		"I",
		"flight",
		"the",
		"morning",
		"through",
		"Denver",
		".",
	}

	tree := Parse(0, sentence)
	t.Log(tree)

	got := make([]string, 0)
	tree.Walk(func(t *tokenize.Token) bool {
		got = append(got, t.Text.Content)
		return true
	})

	if !slices.Equal(got, want) {
		t.Errorf("Tree.Walk() = %v, want %v", got, want)
	}
}

func TestTreeDependencies(t *testing.T) {
	t.Parallel()

	{
		sentence := testhelper.NewExampleTokens1(t, 0, 0)
		tree := Parse(0, sentence)

		want := [][2]string{
			{"prefer", "I"},
			{"prefer", "flight"},
			{"flight", "the"},
			{"flight", "morning"},
			{"flight", "through"},
			{"through", "Denver"},
		}

		got := make([][2]string, 0)
		tree.Dependencies(func(head, dependent *tokenize.Token) bool {
			got = append(got, [2]string{head.Text.Content, dependent.Text.Content})
			return true
		})
		for _, edge := range want {
			if !slices.Contains(want, edge) {
				t.Errorf("Tree.Dependencies() = %v, want %v", got, want)
			}
		}
	}
	{
		sentence := testhelper.NewExampleTokens1(t, 0, 0)
		tree := Parse(0, sentence)

		want := [][2]string{
			{"prefer", "I"},
			{"prefer", "flight"},
		}
		got := make([][2]string, 0)
		tree.Dependencies(func(from, to *tokenize.Token) bool {
			got = append(got, [2]string{from.Text.Content, to.Text.Content})
			return from.Text.Content != "flight"
		})

		for _, edge := range want {
			if !slices.Contains(want, edge) {
				t.Errorf("Tree.Dependencies() = %v, want %v", got, want)
			}
		}
	}
}

func TestTreeHead(t *testing.T) {
	t.Parallel()

	{
		tokens := testhelper.NewExampleTokens1(t, 0, 0)
		tree := Parse(0, tokens)
		tok := tokens[4] // flight
		got := tree.Head(tok)

		want := tokens[1] // prefer
		if got != want {
			t.Errorf("Tree.Head() = %v, want %v", got, want)
		}
	}
}
