package entitydebs

import (
	"encoding/json"

	"github.com/ndabAP/entitydebs/tokenize"
)

func (f Frames) MarshalJSON() ([]byte, error) {
	type frame struct {
		Sentences []*tokenize.Sentence `json:"sentences"`
		Tokens    []*tokenize.Token    `json:"tokens"`
		Sentiment *tokenize.Sentiment  `json:"sentiment"`
	}

	// Serialize frames.
	frames := make([]frame, 0, len(f.frames))
	for _, f := range f.frames {
		frames = append(frames, frame{
			Sentences: f.sentences,
			Tokens:    f.tokens,
			Sentiment: f.sentiment,
		})
	}
	return json.Marshal(struct {
		Frames   []frame                     `json:"frames"`
		Entities map[string][]tokenize.Token `json:"entities"`
	}{
		Frames:   frames,
		Entities: f.entities,
	})
}
