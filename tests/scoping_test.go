// tests/scoping_test.go
package tests

import (
	"context"
	"testing"

	"github.com/smegg99/s99logger"
)

func TestMinLevelFiltering(t *testing.T) {
	sink := &captureSink{}
	logger := s99logger.New(sink, s99logger.Options{MinLevel: s99logger.LevelWarn})

	logger.Debug(event{id: "d"})
	logger.Info(event{id: "i"})
	logger.Warn(event{id: "w"})
	logger.Error(event{id: "e"})

	if len(sink.records) != 2 {
		t.Fatalf("got %d records, want 2 (warn, error)", len(sink.records))
	}
	if sink.records[0].MessageID != "w" || sink.records[1].MessageID != "e" {
		t.Errorf("unexpected records: %+v", sink.records)
	}
}

func TestWithAttrs(t *testing.T) {
	sink := &captureSink{}
	base := s99logger.New(sink, s99logger.Options{})
	scoped := base.With(s99logger.String("request_id", "abc"))

	scoped.Info(event{id: "evt", attrs: []s99logger.Attr{s99logger.Int("n", 1)}})

	got := attrMap(sink.records[0].Attrs)
	if got["request_id"] != "abc" {
		t.Errorf("missing base attr request_id: %+v", sink.records[0].Attrs)
	}
	if got["n"] != 1 {
		t.Errorf("missing event attr n: %+v", sink.records[0].Attrs)
	}

	// The original logger must be unaffected (With returns a child).
	base.Info(event{id: "evt2"})
	if _, leaked := attrMap(sink.records[1].Attrs)["request_id"]; leaked {
		t.Error("base logger leaked a With attr")
	}
}

func TestContextLanguageOverridesOption(t *testing.T) {
	sink := &captureSink{}
	logger := s99logger.New(sink, s99logger.Options{
		Language:   "en",
		Translator: newTranslator(t),
	})
	greeting := event{id: "greeting", attrs: []s99logger.Attr{s99logger.String("name", "Ola")}}

	// A language on the context wins over Options.Language.
	logger.Info(greeting, s99logger.ContextWithLanguage(context.Background(), "pl"))
	if sink.records[0].Message != "Cześć Ola" {
		t.Errorf("message = %q, want Polish from context", sink.records[0].Message)
	}
	if sink.records[0].Lang != "pl" {
		t.Errorf("record lang = %q, want pl", sink.records[0].Lang)
	}

	// Without a context language, it falls back to Options.Language (en).
	logger.Info(greeting)
	if sink.records[1].Message != "Hello Ola" {
		t.Errorf("message = %q, want English default", sink.records[1].Message)
	}
}

func TestLanguageContextUsesBackground(t *testing.T) {
	ctx := s99logger.LanguageContext("pl")

	if got, ok := s99logger.LanguageFromContext(ctx); !ok || got != "pl" {
		t.Fatalf("LanguageFromContext = %q, %t; want pl, true", got, ok)
	}
}

func attrMap(attrs []s99logger.Attr) map[string]any {
	m := make(map[string]any, len(attrs))
	for _, a := range attrs {
		m[a.Key] = a.Value
	}
	return m
}
