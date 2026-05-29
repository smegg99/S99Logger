package i18n

import (
	"io/fs"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

// Decoder unmarshals a locale file of some format into go-i18n messages. It is
// the same signature go-i18n uses, re-exported so callers need not import it.
type Decoder = goi18n.UnmarshalFunc

// Options configures the Translator.
type Options struct {
	// DefaultLanguage is the fallback language tag (e.g. "en"). Required.
	DefaultLanguage string
	// FS holds the locale files. Works with embed.FS and os.DirFS. Required.
	FS fs.FS
	// Files are the locale file paths within FS to load. The language is
	// derived from each file name (e.g. "locales/pl.toml" -> Polish).
	Files []string
	// Decoders registers additional formats by name (e.g. "yaml", "toml").
	// JSON is always registered.
	Decoders map[string]Decoder
}
