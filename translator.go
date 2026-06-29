// translator.go

package s99logger

import "context"

// Translator turns a typed event into a localized message for a language.
// It is implemented outside this package (e.g. by the i18n adapter) so that
// this package never depends on any translation library.
//
// When Translate returns an error, the Logger still emits the record, falling
// back to the event's MessageID as the message.
type Translator interface {
	Translate(ctx context.Context, lang string, event Event) (string, error)
}
