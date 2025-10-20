package testhelper

import (
	"testing"

	"github.com/ndabAP/entitydebs/tokenize"
)

type PoSParams struct {
	Tag         tokenize.PartOfSpeechTag
	Aspect      tokenize.PartOfSpeechAspect
	Case        tokenize.PartOfSpeechCase
	Form        tokenize.PartOfSpeechForm
	Gender      tokenize.PartOfSpeechGender
	Mood        tokenize.PartOfSpeechMood
	Number      tokenize.PartOfSpeechNumber
	Person      tokenize.PartOfSpeechPerson
	Proper      tokenize.PartOfSpeechProper
	Reciprocity tokenize.PartOfSpeechReciprocity
	Tense       tokenize.PartOfSpeechTense
	Voice       tokenize.PartOfSpeechVoice
}

func NewPoS(t *testing.T, tag tokenize.PartOfSpeechTag, params PoSParams) *tokenize.PartOfSpeech {
	t.Helper()

	pos := &tokenize.PartOfSpeech{
		Tag:         tag,
		Aspect:      params.Aspect,
		Case:        params.Case,
		Form:        params.Form,
		Gender:      params.Gender,
		Mood:        params.Mood,
		Number:      params.Number,
		Person:      params.Person,
		Proper:      params.Proper,
		Reciprocity: params.Reciprocity,
		Tense:       params.Tense,
		Voice:       params.Voice,
	}

	return pos
}

func NewDepEdge(t *testing.T,
	headTokenIndex int32,
	label tokenize.DependencyEdgeLabel,
) *tokenize.DependencyEdge {
	t.Helper()

	return &tokenize.DependencyEdge{
		HeadTokenIndex: headTokenIndex,
		Label:          label,
	}
}

func NewToken(t *testing.T,
	content string,
	beginOffset int32,
	pos *tokenize.PartOfSpeech,
	depEdge *tokenize.DependencyEdge,
	lema string,
) *tokenize.Token {
	t.Helper()

	return &tokenize.Token{
		Text: &tokenize.TextSpan{
			Content:     content,
			BeginOffset: beginOffset,
		},
		PartOfSpeech:   pos,
		DependencyEdge: depEdge,
		Lemma:          lema,
	}
}

func NewSentence(t *testing.T,
	content string,
	beginOffset int32,
	sentiment *tokenize.Sentiment,
) *tokenize.Sentence {
	t.Helper()

	return &tokenize.Sentence{
		Text: &tokenize.TextSpan{
			Content:     content,
			BeginOffset: beginOffset,
		},
		Sentiment: sentiment,
	}
}
