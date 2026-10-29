// options.go

package s99logger

import "time"

// Options configures a Logger; the zero value works, with time.Now and untranslated message ids.
type Options struct {
	// Language is the default target language passed to the Translator. A call logged with
	// LanguageContext(lang) overrides it.
	Language string
	// Service is the service name attached to every record.
	Service string
	// MinLevel is the lowest level that is logged; lower levels are dropped
	// before any translation or sink work. The zero value (LevelDebug) logs
	// everything.
	MinLevel Level
	// Translator localizes events. If nil, messages are the message id.
	Translator Translator
	// Clock supplies the record timestamp. If nil, time.Now is used.
	Clock func() time.Time
}
