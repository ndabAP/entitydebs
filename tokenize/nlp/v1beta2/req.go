package v1beta2

import (
	"context"

	apiv1beta2 "cloud.google.com/go/language/apiv1beta2"
	"cloud.google.com/go/language/apiv1beta2/languagepb"
	"github.com/ndabAP/entitydebs/tokenize/nlp/internal/retry"
	"github.com/ndabAP/entitydebs/tokenize/nlp/language"
	"google.golang.org/api/option"
)

type api struct {
	creds string
	lang  string
}

func New(creds, lang string) api {
	return api{
		creds: creds,
		lang:  lang,
	}
}

func (v1 api) Syntax(ctx context.Context, text string) (*languagepb.AnalyzeSyntaxResponse, error) {
	client, err := apiv1beta2.NewClient(
		ctx,
		option.WithCredentialsFile(v1.creds),
	)
	if err != nil {
		return &languagepb.AnalyzeSyntaxResponse{}, err
	}
	//nolint:errcheck
	defer client.Close()

	doc := &languagepb.Document{
		Source: &languagepb.Document_Content{
			Content: text,
		},
		Type: languagepb.Document_PLAIN_TEXT,
	}
	if v1.lang != language.Auto {
		doc.Language = v1.lang
	}

	var res *languagepb.AnalyzeSyntaxResponse
	if err := retry.Do(ctx, func() error {
		res, err = client.AnalyzeSyntax(ctx, &languagepb.AnalyzeSyntaxRequest{
			Document:     doc,
			EncodingType: languagepb.EncodingType_UTF8,
		})

		return err
	}); err != nil {
		return res, err
	}

	return res, err
}
