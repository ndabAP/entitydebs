package dependency

import (
	"testing"

	"github.com/ndabAP/entitydebs/testhelper"
	"github.com/ndabAP/entitydebs/tokenize"
)

func TestParse(t *testing.T) {
	{
		var (
			sentence1 = testhelper.NewExampleTokens1(t, 0, 0)

			charOffset = sentence1[len(sentence1)-1].Text.BeginOffset + 2 // Terminal plus white space.
			depOffset  = int32(len(sentence1)) - 1
			sentence2  = testhelper.NewExampleTokens2(t, charOffset, depOffset)
		)

		tree := Parse(depOffset, sentence2)
		t.Log(tree)

		edge := tree.graph.Edge(sentence2[0].ID(), sentence2[1].ID())
		t.Log(edge)
		var (
			got  = edge.From().ID()
			want = sentence2[0].ID()
		)
		if got != want {
			t.Errorf("Tree.graph.Edge() = %v, want %v", got, want)
		}
	}
	{
		toks := []*tokenize.Token{
			{
				Text:           &tokenize.TextSpan{},
				DependencyEdge: testhelper.NewDepEdge(t, 1, tokenize.DependencyEdgeLabelRoot),
			},
			{
				Text:           &tokenize.TextSpan{BeginOffset: 1},
				DependencyEdge: testhelper.NewDepEdge(t, 0, tokenize.DependencyEdgeLabelUnknown),
			},
		}
		defer func() { _ = recover() }()
		Parse(0, toks)
		t.Error("Parse() expected to panic on cyclic dependency")
	}
}
