package tokenize

type Token struct {
	Text           *TextSpan
	PartOfSpeech   *PartOfSpeech
	DependencyEdge *DependencyEdge
	Lemma          string
}

// NilToken can be used as a placeholder token to not lose positional
// properties. Note that the NilToken loses dependency relations.
var NilToken = &Token{
	Text: &TextSpan{
		BeginOffset: -1,
	},
	PartOfSpeech: &PartOfSpeech{
		Tag:         -1,
		Aspect:      -1,
		Case:        -1,
		Form:        -1,
		Gender:      -1,
		Mood:        -1,
		Number:      -1,
		Person:      -1,
		Proper:      -1,
		Reciprocity: -1,
		Tense:       -1,
		Voice:       -1,
	},
	DependencyEdge: &DependencyEdge{
		HeadTokenIndex: -1,
		Label:          -1,
	},
}

func (t Token) ID() int64 {
	return int64(t.Text.BeginOffset)
}

func (t Token) String() string {
	return t.Text.Content
}

func (t Token) Clone() (token *Token) {
	token = &Token{}
	if t.Text != nil {
		token.Text = &TextSpan{}
		token.Text.BeginOffset = t.Text.BeginOffset
		token.Text.Content = t.Text.Content
	}
	if t.PartOfSpeech != nil {
		token.PartOfSpeech = &PartOfSpeech{}
		token.PartOfSpeech.Tag = t.PartOfSpeech.Tag
		token.PartOfSpeech.Aspect = t.PartOfSpeech.Aspect
		token.PartOfSpeech.Case = t.PartOfSpeech.Case
		token.PartOfSpeech.Form = t.PartOfSpeech.Form
		token.PartOfSpeech.Gender = t.PartOfSpeech.Gender
		token.PartOfSpeech.Mood = t.PartOfSpeech.Mood
		token.PartOfSpeech.Number = t.PartOfSpeech.Number
		token.PartOfSpeech.Person = t.PartOfSpeech.Person
		token.PartOfSpeech.Proper = t.PartOfSpeech.Proper
		token.PartOfSpeech.Reciprocity = t.PartOfSpeech.Reciprocity
		token.PartOfSpeech.Tense = t.PartOfSpeech.Tense
		token.PartOfSpeech.Voice = t.PartOfSpeech.Voice
	}
	if t.DependencyEdge != nil {
		token.DependencyEdge = &DependencyEdge{}
		token.DependencyEdge.HeadTokenIndex = t.DependencyEdge.HeadTokenIndex
		token.DependencyEdge.Label = t.DependencyEdge.Label
	}
	token.Lemma = t.Lemma
	return
}

func (t Token) DOTID() string {
	return t.Text.Content
}
