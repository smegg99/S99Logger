package s99logger

import "time"

// Options configures a Logger. The zero value is usable: Clock defaults to
// time.Now and Translator may be nil (messages fall back to the message id).
type Options struct {
	// Language is the default target language passed to the Translator. It can
	// be overridden per call with ContextWithLanguage.
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
