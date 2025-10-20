package tokenize

import (
	"fmt"
	"strings"
)

// Analysis contains the sentences, tokens and sentiment of a tokenized text.
type Analysis struct {
	// Sentences contains each sentence's text and sentiment.
	Sentences []*Sentence
	// Tokens contains all document tokens.
	Tokens []*Token
	// Sentiment is the documents Sentiment.
	Sentiment *Sentiment
}

func (a Analysis) String() string {
	var sb strings.Builder

	// Sentences
	for i, sentence := range a.Sentences {
		fmt.Fprintf(&sb, "index:%d ", i)
		sb.WriteString("\n")

		if sentence.Text != nil {
			fmt.Fprintf(&sb, "text: content:%s begin_offset:%d",
				sentence.Text.Content,
				sentence.Text.BeginOffset,
			)
			sb.WriteString("\n")
		}
		if sentence.Sentiment != nil {
			fmt.Fprintf(&sb, "sentiment: magnitude:%f score:%f",
				sentence.Sentiment.Magnitude,
				sentence.Sentiment.Score,
			)
		}

		sb.WriteString("\n")
	}

	// Tokens
	for i, token := range a.Tokens {
		fmt.Fprintf(&sb, "index:%d", i)
		sb.WriteString("\n")

		if token.Text != nil {
			fmt.Fprintf(&sb, "text: content:%s begin_offset:%d",
				token.Text.Content,
				token.Text.BeginOffset,
			)
			sb.WriteString("\n")
		}
		if token.PartOfSpeech != nil {
			fmt.Fprintf(&sb, "part_of_speech: tag:%d aspect:%d case:%d form:%d gender:%d mood:%d number:%d person:%d proper:%d reciprocity:%d tense:%d voice:%d",
				token.PartOfSpeech.Tag,
				token.PartOfSpeech.Aspect,
				token.PartOfSpeech.Case,
				token.PartOfSpeech.Form,
				token.PartOfSpeech.Gender,
				token.PartOfSpeech.Mood,
				token.PartOfSpeech.Number,
				token.PartOfSpeech.Person,
				token.PartOfSpeech.Proper,
				token.PartOfSpeech.Reciprocity,
				token.PartOfSpeech.Tense,
				token.PartOfSpeech.Voice,
			)
			sb.WriteString("\n")
		}
		if token.DependencyEdge != nil {
			fmt.Fprintf(&sb, "dependency_edge: head_token_index:%d label:%d",
				token.DependencyEdge.HeadTokenIndex,
				token.DependencyEdge.Label,
			)
			sb.WriteString("\n")
		}
		fmt.Fprintf(&sb, "lemma:%s\n", token.Lemma)

		sb.WriteString("\n")
	}

	// Sentiment
	if a.Sentiment != nil {
		fmt.Fprintf(&sb, "sentiment: magnitude:%f score:%f\n",
			a.Sentiment.Magnitude,
			a.Sentiment.Score,
		)

		sb.WriteString("\n")
	}

	fmt.Fprintf(&sb, "%#v", a)

	sb.WriteString("\n")

	return sb.String()
}
