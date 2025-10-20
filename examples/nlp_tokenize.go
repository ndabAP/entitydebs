//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"

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
		"I prefer the morning flight through Denver.",
		"The quick brown fox jumps over the lazy dog's back",
	}
	nlp := nlp.New(creds, language.Auto)
	for i, text := range texts {
		analysis, err := nlp.Tokenize(ctx, text, tokenize.FeatureAll)
		if err != nil {
			panic(err.Error())
		}
		fmt.Fprintln(os.Stdout, analysis.String())
		if i != len(texts)-1 {
			fmt.Println()
		}
	}
}
