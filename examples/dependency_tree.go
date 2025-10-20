//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"cloud.google.com/go/language/apiv1/languagepb"
	"github.com/ndabAP/entitydebs"
	"github.com/ndabAP/entitydebs/dependency"
	"github.com/ndabAP/entitydebs/tokenize"
	"github.com/ndabAP/entitydebs/tokenize/nlp"
	"github.com/ndabAP/entitydebs/tokenize/nlp/language"
)

func main() {
	var (
		ctx = context.Background()

		creds = os.Getenv("GCLOUD_SERVICE_ACCOUNT_KEY")
	)

	texts := []string{
		"Ask not what your country can do for you - ask what you can do for your country.",
	}
	entity := []string{"country"}
	src := entitydebs.NewSource(
		entity,
		texts,
	)
	nlp := nlp.New(creds, language.Auto)
	frames, err := src.Frames(
		ctx,
		nlp,
		tokenize.FeatureSyntax,
		entitydebs.Lowercaser,
	)
	if err != nil {
		panic(err.Error())
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "Head\tRelationship\tDependent")
	frames.Forest().Dependencies(func(head, dependent *tokenize.Token, rel tokenize.DependencyEdgeLabel, tree dependency.Tree) bool {
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\n",
			head.Text.Content,
			languagepb.DependencyEdge_Label_name[int32(rel)],
			dependent.Text.Content,
		)
		return true
	})
	w.Flush()
	fmt.Println()

	heads := frames.Forest().Heads(nil)
	toks := make([]string, len(heads))
	for i, head := range heads {
		toks[i] = head.String()
	}
	fmt.Fprintf(os.Stdout, "heads: %s", strings.Join(toks, ","))

	dependents := frames.Forest().Dependents(nil)
	toks = make([]string, len(dependents))
	for i, dependent := range dependents {
		toks[i] = dependent.String()
	}
	fmt.Fprintf(os.Stdout, "\ndependents: %s", strings.Join(toks, ","))
}
