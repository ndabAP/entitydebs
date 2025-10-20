//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ndabAP/entitydebs"
	"github.com/ndabAP/entitydebs/tokenize"
	"github.com/ndabAP/entitydebs/tokenize/nlp"
	"github.com/ndabAP/entitydebs/tokenize/nlp/language"
)

func main() {
	var (
		ctx = context.Background()

		creds = os.Getenv("GCLOUD_SERVICE_ACCOUNT_KEY")
	)

	src := entitydebs.NewSource(
		[]string{"Max Payne", "Payne"},
		[]string{
			"Yeah, somethings wrong with the bosses. Payne's there, and they're not answering.",
			"Bang! You're dead, Max Payne!",
		},
	)
	nlp := nlp.New(creds, language.EN)
	frames, err := src.Frames(
		ctx,
		nlp,
		tokenize.FeatureAll,
		entitydebs.Lemma,
		entitydebs.NFKC,
		entitydebs.Lowercaser,
	)
	if err != nil {
		panic(err.Error())
	}

	for _, tokens := range frames.All() {
		for _, token := range tokens {
			fmt.Fprintln(os.Stdout, token.String())
		}
	}
}
