# entitydebs

[![Go Report Card](https://goreportcard.com/badge/github.com/ndabAP/entitydebs)](https://goreportcard.com/report/github.com/ndabAP/entitydebs)
[![go.dev reference](https://pkg.go.dev/badge/github.com/ndabAP/entitydebs)](https://pkg.go.dev//github.com/ndabAP/entitydebs)
![GitHub Tag](https://img.shields.io/github/v/tag/ndabAP/entitydebs)
[![Build Status](https://github.com/ndabAP/entitydebs/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/ndabAP/entitydebs/actions/workflows/test.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ndabAP/entitydebs)
![GitHub License](https://img.shields.io/github/license/ndabAP/entitydebs)

`entitydebs` is a social science tool to programmatically analyze entities in
non-fictional texts. In particular, it's well-suited to extract the sentiment
for an entity using dependency parsing. Tokenization is highly customizable and
supports the Google Cloud Natural Language API out-of-the-box. It can help
answer questions like:

- How do politicians describe their country in governmental speeches?
- Which current topics correlate with celebrities?
- What are the most common root verbs used in different music genres?

Visit the [live demo](https://ndabap.github.io/entityscrape/) or read through
the source code [here](https://github.com/ndabAP/entityscrape). To learn more
about dependency trees consult the [Google Natural Language API guide](https://cloud.google.com/natural-language/docs/morphology#dependency_trees).


## Features

- **Dependency parsing**: Build and traverse dependency trees for syntactic and
sentiment analysis
- **AI tokenizer**: Out-of-the-box support for the [Google Cloud Natural
Language API](https://cloud.google.com/natural-language?hl=en) for robust
tokenization, with a built-in retrier
- **Bullet-proof trees**: Dependency trees are constructed using
[gonum](https://github.com/gonum/gonum)
- **Efficient traversal**: Native iterators for traversing analysis results
- **Text normalization**: Built-in normalizers (lowercasing, NFKC,
lemmatization) to reduce redundancy and improve data integrity
- **High test coverage**: Over 80 % test coverage and millions of tokens

## Install

```sh
go get github.com/ndabAP/entitydebs
```

## Example

Let's understand how to use the tool with an example. We would like to know
the dependencies in Congress speeches to the United States. Our entity would be
United States and our texts Congressional speeches.

First, we need to create a `source` to create the data frames from. Our entity
with its aliases is:

```go
entity := []string{"America", "USA", "United States", "US"}
```

As text source we can use ["Congressional Record for the 43rd-114th Congresses: Parsed Speeches and Phrase Counts"](https://data.stanford.edu/congress_text).
We skip the part how to parse the texts and use a subset for this example:

```go
texts := []string{
	"Pay tribute for the past 237 years of sacrifice to our great United States Army.",
	"So having a Special Order this evening is an opportunity for us all to come together and celebrate the commitment of the United States Congress to communities around the world as they experience.",
	"Thank you very much for your sacrifice and your commitment for a free Cuba and a strong United States.",
}
src := entitydebs.NewSource(entity, texts)
```

Now, we need a tokenizer that implements `tokenize.Tokenizer`. The built-in
package `nlp` uses Google Natural Language API and support dependency parsing.
For this to work, we need a service account file
(learn more [here](https://cloud.google.com/iam/docs/service-account-overview))
with the respective permissions.

```go
creds := os.Getenv("GCLOUD_SERVICE_ACCOUNT_KEY")
nlp := nlp.New(creds, language.EN)
```

`source.Frames` uses the provided tokenizer to generate the data frames. This
may take a while depending on the input and how the tokenizer works.

```go
frames, err := src.Frames(ctx, nlp, tokenize.FeatureSyntax)
if err != nil {
	panic(err.Error())
}
```

From this point, `Frames` can construct for us a forest of dependency trees.
Trees that don't contain the entity are not created. We print all relationships
using Gos native `text/tabwriter`:

```go
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
```

You can find more examples in the
[examples](https://github.com/ndabAP/entitydebs/tree/main/examples) folder.

## Author

[Julian Claus](https://www.julian-claus.de) and contributors.

## License

MIT
