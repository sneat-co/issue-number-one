package translations

import (
	"cloud.google.com/go/translate"
	"context"
	"errors"
	"fmt"
	"golang.org/x/text/language"
)

type GoogleTranslator struct{ Client *translate.Client }

// NewGoogleTranslator uses Google Application Default Credentials in the host project.
func NewGoogleTranslator(ctx context.Context) (*GoogleTranslator, error) {
	client, err := translate.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("create Google Translation client: %w", err)
	}
	return &GoogleTranslator{Client: client}, nil
}

func (translator *GoogleTranslator) Close() error {
	if translator == nil || translator.Client == nil {
		return nil
	}
	return translator.Client.Close()
}

func (translator *GoogleTranslator) Translate(ctx context.Context, sourceLanguage, targetLanguage string, texts []string) ([]string, error) {
	if translator == nil || translator.Client == nil {
		return nil, errors.New("Google Translation client is required")
	}
	source, err := language.Parse(sourceLanguage)
	if err != nil {
		return nil, fmt.Errorf("parse source language: %w", err)
	}
	target, err := language.Parse(targetLanguage)
	if err != nil {
		return nil, fmt.Errorf("parse target language: %w", err)
	}
	response, err := translator.Client.Translate(ctx, texts, target, &translate.Options{Source: source})
	if err != nil {
		return nil, err
	}
	result := make([]string, len(response))
	for index := range response {
		result[index] = response[index].Text
	}
	return result, nil
}
