//go:build ignore

package entitydebs_test

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"cloud.google.com/go/language/apiv1/languagepb"
	"github.com/ndab!a!p/entitydebs/v14"
	"github.com/ndabAP/entitydebs/dependency"
	"github.com/ndabAP/entitydebs/tokenize"
	"github.com/ndabAP/entitydebs/tokenize/nlp"
	"github.com/ndabAP/entitydebs/tokenize/nlp/language"
)

func ExampleNewSource() {
	ctx := context.Background()

	// Create a new source. You need to parse the texts into a human-readable
	// format before.
	entity := []string{"America", "USA", "United States", "US"}
	texts := []string{
		"Pay tribute for the past 237 years of sacrifice to our great United States Army.",
		"So having a Special Order this evening is an opportunity for us all to come together and celebrate the commitment of the United States Congress to communities around the world as they experience.",
		"Thank you very much for your sacrifice and your commitment for a free Cuba and a strong United States.",
	}
	src := entitydebs.NewSource(entity, texts)

	// Create a tokenizer. In this case, we use the built-in Google Natural
	// Language API.
	creds := os.Getenv("GCLOUD_SERVICE_ACCOUNT_KEY")
	nlp := nlp.New(creds, language.EN)

	// Now we can create a data frames and perfom the API request. This may take
	// a while.
	frames, err := src.Frames(ctx, nlp, tokenize.FeatureSyntax)
	if err != nil {
		panic(err.Error())
	}

	// Print all dependency relationships that contain entity tokens.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	_, _ = fmt.Fprintln(w, "Head\tRelationship\tDependent")
	frames.Forest().Dependencies(func(
		head,
		dependent *tokenize.Token,
		rel tokenize.DependencyEdgeLabel,
		tree dependency.Tree,
	) bool {
		_, _ = fmt.Fprintf(
			w,
			"%s\t%s\t%s\n",
			head.Text.Content,
			languagepb.DependencyEdge_Label_name[int32(rel)],
			dependent.Text.Content,
		)
		return true
	})
	_ = w.Flush()
	// Output:
	// Head        Relationship Dependent
	// tribute     P            .
	// tribute     PREP         to
	// tribute     PREP         for
	// tribute     NN           Pay
	// for         POBJ         years
	// years       PREP         of
	// years       NUM          237
	// years       AMOD         past
	// years       DET          the
	// of          POBJ         sacrifice
	// to          POBJ         Army
	// States      NN           United
	// Army        NN           States
	// Army        AMOD         great
	// Army        POSS         our
	// having      TMOD         evening
	// having      DOBJ         Order
	// having      ADVMOD       So
	// Order       NN           Special
	// Order       DET          a
	// evening     DET          this
	// is          P            .
	// is          ATTR         opportunity
	// is          CSUBJ        having
	// opportunity CCOMP        come
	// opportunity DET          an
	// us          DET          all
	// come        CONJ         celebrate
	// come        CC           and
	// come        ADVMOD       together
	// come        AUX          to
	// come        NSUBJ        us
	// come        MARK         for
	// celebrate   ADVCL        experience
	// celebrate   PREP         to
	// celebrate   DOBJ         commitment
	// commitment  PREP         of
	// commitment  DET          the
	// of          POBJ         Congress
	// States      NN           United
	// Congress    NN           States
	// Congress    DET          the
	// to          POBJ         communities
	// communities PREP         around
	// around      POBJ         world
	// world       DET          the
	// experience  NSUBJ        they
	// experience  MARK         as
	// Thank       P            .
	// Thank       PREP         for
	// Thank       ADVMOD       much
	// Thank       DOBJ         you
	// much        ADVMOD       very
	// for         POBJ         sacrifice
	// sacrifice   CONJ         commitment
	// sacrifice   CC           and
	// sacrifice   POSS         your
	// commitment  PREP         for
	// commitment  POSS         your
	// for         POBJ         Cuba
	// Cuba        CONJ         States
	// Cuba        CC           and
	// Cuba        AMOD         free
	// Cuba        DET          a
	// States      NN           United
	// States      AMOD         strong
	// States      DET          a
}
