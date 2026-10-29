// event.go

package s99logger

// MessageID is a stable, language-independent message id, emitted beside the text and used as the text when translation fails.
type MessageID string

// Event is a typed log event; NewEvent and NewPluralEvent cover the simple cases.
type Event interface {
	// MessageID returns the stable identifier for this event.
	MessageID() MessageID
	// Attrs returns the structured attributes for this event.
	Attrs() []Attr
	// PluralCount returns the count used for plural translation selection.
	// The boolean is false for events that are not pluralized.
	PluralCount() (int, bool)
}

// NoPlural lets an event struct satisfy Event without a PluralCount of its own.
type NoPlural struct{}

// PluralCount reports that the event is not pluralized.
func (NoPlural) PluralCount() (int, bool) { return 0, false }
