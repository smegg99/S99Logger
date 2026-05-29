// tests/i18n_test.go
package tests

import (
	"context"
	"embed"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/smegg99/s99logger"
	"github.com/smegg99/s99logger/i18n"
)

//go:embed testdata/en.json testdata/pl.json testdata/en.toml
var localesFS embed.FS

func newTranslator(t *testing.T) *i18n.Translator {
	t.Helper()
	tr, err := i18n.New(i18n.Options{
		DefaultLanguage: "en",
		FS:              localesFS,
		Files:           []string{"testdata/en.json", "testdata/pl.json"},
	})
	if err != nil {
		t.Fatalf("i18n.New: %v", err)
	}
	return tr
}

func TestI18nLoadsEmbeddedJSON(t *testing.T) {
	tr := newTranslator(t)

	got, err := tr.Translate(context.Background(), "pl", event{
		id:    "greeting",
		attrs: []s99logger.Attr{s99logger.String("name", "Karol")},
	})
	if err != nil {
		t.Fatalf("Translate: %v", err)
	}
	if got != "Cześć Karol" {
		t.Errorf("got %q, want %q", got, "Cześć Karol")
	}
}

func TestI18nLoadsEmbeddedTOML(t *testing.T) {
	tr, err := i18n.New(i18n.Options{
		DefaultLanguage: "en",
		FS:              localesFS,
		Files:           []string{"testdata/en.toml"},
		Decoders:        map[string]i18n.Decoder{"toml": toml.Unmarshal},
	})
	if err != nil {
		t.Fatalf("i18n.New: %v", err)
	}

	got, err := tr.Translate(context.Background(), "en", event{
		id:    "greeting",
		attrs: []s99logger.Attr{s99logger.String("name", "Karol")},
	})
	if err != nil {
		t.Fatalf("Translate: %v", err)
	}
	if got != "Hello Karol from TOML" {
		t.Errorf("got %q, want %q", got, "Hello Karol from TOML")
	}
}

func TestI18nFallbackLanguage(t *testing.T) {
	tr := newTranslator(t)

	// "fr" has no locale file, so it must fall back to the default (en).
	got, err := tr.Translate(context.Background(), "fr", event{
		id:    "greeting",
		attrs: []s99logger.Attr{s99logger.String("name", "Karol")},
	})
	if err != nil {
		t.Fatalf("Translate: %v", err)
	}
	if got != "Hello Karol" {
		t.Errorf("got %q, want fallback %q", got, "Hello Karol")
	}
}

func TestI18nPluralization(t *testing.T) {
	tr := newTranslator(t)
	ctx := context.Background()

	cases := []struct {
		lang  string
		count int
		want  string
	}{
		{"en", 1, "1 item"},
		{"en", 3, "3 items"},
		{"pl", 1, "1 element"},
		{"pl", 3, "3 elementy"},
		{"pl", 5, "5 elementów"},
	}
	for _, c := range cases {
		got, err := tr.Translate(ctx, c.lang, event{
			id:     "items",
			attrs:  []s99logger.Attr{s99logger.Int("count", c.count)},
			count:  c.count,
			plural: true,
		})
		if err != nil {
			t.Fatalf("Translate(%s,%d): %v", c.lang, c.count, err)
		}
		if got != c.want {
			t.Errorf("Translate(%s,%d) = %q, want %q", c.lang, c.count, got, c.want)
		}
	}
}

func TestI18nAttrsPassedIntoTemplates(t *testing.T) {
	tr := newTranslator(t)

	got, err := tr.Translate(context.Background(), "en", event{
		id:    "greeting",
		attrs: []s99logger.Attr{s99logger.String("name", "World")},
	})
	if err != nil {
		t.Fatalf("Translate: %v", err)
	}
	if got != "Hello World" {
		t.Errorf("attr not rendered into template: got %q", got)
	}
}

func TestI18nVerifyCoverage(t *testing.T) {
	tr := newTranslator(t)
	if err := tr.Verify([]string{"en", "pl"}, "greeting", "items"); err != nil {
		t.Errorf("Verify reported missing translations for present ids: %v", err)
	}
}

func TestI18nVerifyReportsMissing(t *testing.T) {
	tr := newTranslator(t)
	err := tr.Verify([]string{"en", "pl"}, "greeting", "does_not_exist")
	if err == nil {
		t.Fatal("expected an error for a missing message id")
	}
	if !strings.Contains(err.Error(), "does_not_exist") {
		t.Errorf("error should name the missing id, got: %v", err)
	}
}

func TestI18nMissingFSReturnsError(t *testing.T) {
	if _, err := i18n.New(i18n.Options{DefaultLanguage: "en"}); err == nil {
		t.Fatal("expected error when FS is nil")
	}
}

func TestI18nInvalidDefaultLanguageReturnsError(t *testing.T) {
	if _, err := i18n.New(i18n.Options{DefaultLanguage: "", FS: localesFS}); err == nil {
		t.Fatal("expected error for invalid default language")
	}
}
