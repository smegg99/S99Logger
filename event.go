// event.go

package s99logger

// MessageID is a stable, language-independent identifier for a log message.
// It is always emitted alongside the localized message, and is used as the
// fallback message when no translator is configured or translation fails.
type MessageID string

// Event is a typed log event. Most simple events can use NewEvent or
// NewPluralEvent; applications can still define structs implementing this
// interface when they want stronger domain types or custom behavior.
type Event interface {
	// MessageID returns the stable identifier for this event.
	MessageID() MessageID
	// Attrs returns the structured attributes for this event.
	Attrs() []Attr
	// PluralCount returns the count used for plural translation selection.
	// The boolean is false for events that are not pluralized.
	PluralCount() (int, bool)
}

// NoPlural can be embedded in event structs that are not pluralized, so they
// satisfy Event without implementing PluralCount themselves.
type NoPlural struct{}

// PluralCount reports that the event is not pluralized.
func (NoPlural) PluralCount() (int, bool) { return 0, false }
