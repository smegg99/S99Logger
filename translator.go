// translator.go

package s99logger

import "context"

// Translator localizes an event for a language; on error the Logger falls back to the MessageID.
type Translator interface {
	Translate(ctx context.Context, lang string, event Event) (string, error)
}
