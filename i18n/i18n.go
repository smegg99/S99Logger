// i18n/i18n.go

// Package i18n is the official go-i18n adapter for s99logger. It loads locale
// files from any fs.FS (including embed.FS) into a go-i18n bundle and
// implements s99logger.Translator, using event attrs as template data and the
// event plural count for plural selection.
package i18n

import (
	"context"
	"encoding/json"
	"fmt"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/smegg99/s99logger"
)

// Translator loads locale files into a go-i18n bundle and implements
// s99logger.Translator. It holds no global state.
type Translator struct {
	bundle      *goi18n.Bundle
	defaultLang string
}

var _ s99logger.Translator = (*Translator)(nil)

// New builds a Translator from opts, loading every configured locale file.
func New(opts Options) (*Translator, error) {
	tag, err := language.Parse(opts.DefaultLanguage)
	if err != nil {
		return nil, fmt.Errorf("parse default language %q: %w", opts.DefaultLanguage, err)
	}
	if opts.FS == nil {
		return nil, errFilesystemRequired
	}

	bundle := goi18n.NewBundle(tag)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	for format, fn := range opts.Decoders {
		bundle.RegisterUnmarshalFunc(format, fn)
	}

	for _, file := range opts.Files {
		if _, err := bundle.LoadMessageFileFS(opts.FS, file); err != nil {
			return nil, fmt.Errorf("load locale %q: %w", file, err)
		}
	}

	return &Translator{bundle: bundle, defaultLang: opts.DefaultLanguage}, nil
}

// Translate localizes event for lang, falling back to the default language.
// Event attrs are passed as template data and the plural count, when present,
// drives plural selection.
func (t *Translator) Translate(_ context.Context, lang string, event s99logger.Event) (string, error) {
	localizer := goi18n.NewLocalizer(t.bundle, lang, t.defaultLang)

	attrs := event.Attrs()
	data := make(map[string]any, len(attrs))
	for _, a := range attrs {
		data[a.Key] = a.Value
	}

	cfg := &goi18n.LocalizeConfig{
		MessageID:    string(event.MessageID()),
		TemplateData: data,
	}
	if count, ok := event.PluralCount(); ok {
		cfg.PluralCount = count
	}

	msg, err := localizer.Localize(cfg)
	if err != nil {
		return "", fmt.Errorf("translate %q (lang %q): %w", event.MessageID(), lang, err)
	}
	return msg, nil
}
