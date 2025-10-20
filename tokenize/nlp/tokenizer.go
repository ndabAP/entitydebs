package nlp

import (
	"context"

	"github.com/ndabAP/entitydebs/tokenize"
	"github.com/ndabAP/entitydebs/tokenize/nlp/language"
	"github.com/ndabAP/entitydebs/tokenize/nlp/v1beta2"
	v2 "github.com/ndabAP/entitydebs/tokenize/nlp/v2"
	"golang.org/x/sync/errgroup"
)

// nlp tokenizes a text using Googles Natural Language AI.
type nlp struct {
	creds string
	lang  string
}

// New returns a new Google Natural Language AI tokenizer instance. NLP has a
// built-in retrier.
//
// lang can be either ISO-639-1 or BCP-47 and defaults to [language.Auto] if
// empty.
func New(creds, lang string) tokenize.Tokenizer {
	if creds == "" {
		panic("credentials must not be empty")
	}
	if lang == "" {
		lang = language.Auto
	}
	return nlp{
		creds: creds,
		lang:  lang,
	}
}

// Tokenize implements the [tokenize.Tokenizer] interface.
func (nlp nlp) Tokenize(ctx context.Context, text string, feats tokenize.Features) (tokenize.Analysis, error) {
	var analysis tokenize.Analysis

	var (
		sentences []*tokenize.Sentence
		tokens    []*tokenize.Token
		sentiment = &tokenize.Sentiment{}
	)

	fns := make([]func() error, 0)

	// Analyse syntax
	syntaxfn := func() error {
		res, err := v1beta2.New(nlp.creds, nlp.lang).Syntax(ctx, text)
		if err != nil {
			return err
		}

		sentences = make([]*tokenize.Sentence, len(res.GetSentences()))
		for i, s := range res.GetSentences() {
			sentence := tokenize.Sentence{}
			if s.Text != nil {
				sentence.Text = &tokenize.TextSpan{
					Content:     s.Text.Content,
					BeginOffset: s.Text.BeginOffset,
				}
			}
			if s.Sentiment != nil {
				sentence.Sentiment = &tokenize.Sentiment{
					Magnitude: s.Sentiment.Magnitude,
					Score:     s.Sentiment.Score,
				}
			}
			sentences[i] = &sentence
		}

		tokens = make([]*tokenize.Token, len(res.GetTokens()))
		for i, t := range res.GetTokens() {
			token := tokenize.Token{}
			if t.Text != nil {
				token.Text = &tokenize.TextSpan{
					Content:     t.Text.Content,
					BeginOffset: t.Text.BeginOffset,
				}
			}
			if t.PartOfSpeech != nil {
				token.PartOfSpeech = &tokenize.PartOfSpeech{
					Tag:         (tokenize.PartOfSpeechTag)(t.PartOfSpeech.Tag),
					Aspect:      (tokenize.PartOfSpeechAspect)(t.PartOfSpeech.Aspect),
					Case:        (tokenize.PartOfSpeechCase)(t.PartOfSpeech.Case),
					Form:        (tokenize.PartOfSpeechForm)(t.PartOfSpeech.Form),
					Gender:      (tokenize.PartOfSpeechGender)(t.PartOfSpeech.Gender),
					Mood:        (tokenize.PartOfSpeechMood)(t.PartOfSpeech.Mood),
					Number:      (tokenize.PartOfSpeechNumber)(t.PartOfSpeech.Number),
					Person:      (tokenize.PartOfSpeechPerson)(t.PartOfSpeech.Person),
					Proper:      (tokenize.PartOfSpeechProper)(t.PartOfSpeech.Proper),
					Reciprocity: (tokenize.PartOfSpeechReciprocity)(t.PartOfSpeech.Reciprocity),
					Tense:       (tokenize.PartOfSpeechTense)(t.PartOfSpeech.Tense),
					Voice:       (tokenize.PartOfSpeechVoice)(t.PartOfSpeech.Voice),
				}
			}
			if t.DependencyEdge != nil {
				token.DependencyEdge = &tokenize.DependencyEdge{
					HeadTokenIndex: t.DependencyEdge.HeadTokenIndex,
					Label:          (tokenize.DependencyEdgeLabel)(t.DependencyEdge.Label),
				}
			}
			token.Lemma = t.Lemma
			tokens[i] = &token
		}

		return nil
	}

	// Analyse sentiment
	annotatefn := func(feats v2.Features) func() error {
		return func() error {
			res, err := v2.New(nlp.creds, nlp.lang).Annotate(ctx, text, feats)
			if err != nil {
				return err
			}

			if s := res.GetDocumentSentiment(); s != nil {
				sentiment.Magnitude = s.Magnitude
				sentiment.Score = s.Score
			}
			return nil
		}
	}

	// Syntax
	if feats&tokenize.FeatureSyntax != 0 {
		fns = append(fns, syntaxfn)
	}
	var v2feats v2.Features
	// Sentiment
	if feats&tokenize.FeatureSentiment != 0 {
		v2feats += v2.ExtractSentiment
		fns = append(fns, annotatefn(v2feats))
	}
	// All features
	if feats == tokenize.FeatureAll {
		fns = []func() error{syntaxfn, annotatefn(v2feats)}
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, fn := range fns {
		g.Go(fn)
	}
	if err := g.Wait(); err != nil {
		return analysis, err
	}

	analysis.Tokens = tokens
	analysis.Sentences = sentences
	analysis.Sentiment = sentiment

	return analysis, nil
}
