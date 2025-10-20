package v2

import (
	"context"

	apiv2 "cloud.google.com/go/language/apiv2"
	"cloud.google.com/go/language/apiv2/languagepb"
	"github.com/ndabAP/entitydebs/tokenize/nlp/internal/retry"
	"github.com/ndabAP/entitydebs/tokenize/nlp/language"
	"google.golang.org/api/option"
)

type api struct {
	creds string
	lang  string
}

type Features int

const (
	ExtractSentiment Features = 1 << iota
)

func New(creds, lang string) api {
	return api{
		creds: creds,
		lang:  lang,
	}
}

func (v2 api) Annotate(ctx context.Context, text string, feats Features) (*languagepb.AnnotateTextResponse, error) {
	client, err := apiv2.NewClient(ctx, option.WithCredentialsFile(v2.creds))
	if err != nil {
		return &languagepb.AnnotateTextResponse{}, err
	}
	//nolint:errcheck
	defer client.Close()

	doc := &languagepb.Document{
		Source: &languagepb.Document_Content{
			Content: text,
		},
		Type: languagepb.Document_PLAIN_TEXT,
	}
	if v2.lang != language.Auto {
		doc.LanguageCode = v2.lang
	}

	var res *languagepb.AnnotateTextResponse
	if err := retry.Do(ctx, func() error {
		f := &languagepb.AnnotateTextRequest_Features{}
		if feats&ExtractSentiment != 0 {
			f.ExtractDocumentSentiment = true
		}

		res, err = client.AnnotateText(ctx, &languagepb.AnnotateTextRequest{
			Document: doc,
			Features: f,
		})

		return err
	}); err != nil {
		return res, err
	}

	return res, err
}
