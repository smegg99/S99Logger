// simple_event.go

package s99logger

// SimpleEvent is a small immutable Event implementation for common log events.
// Use NewEvent or NewPluralEvent when a dedicated event struct would only
// repeat the message id and attrs.
type SimpleEvent struct {
	id          MessageID
	attrs       []Attr
	pluralCount int
	plural      bool
}

var _ Event = SimpleEvent{}

// NewEvent returns a non-pluralized event with id and attrs.
func NewEvent(id MessageID, attrs ...Attr) SimpleEvent {
	return SimpleEvent{id: id, attrs: cloneAttrs(attrs)}
}

// NewPluralEvent returns an event whose count drives plural translation
// selection.
func NewPluralEvent(id MessageID, count int, attrs ...Attr) SimpleEvent {
	return SimpleEvent{
		id:          id,
		attrs:       cloneAttrs(attrs),
		pluralCount: count,
		plural:      true,
	}
}

// MessageID returns the stable identifier for this event.
func (e SimpleEvent) MessageID() MessageID { return e.id }

// Attrs returns the structured attributes for this event.
func (e SimpleEvent) Attrs() []Attr { return cloneAttrs(e.attrs) }

// PluralCount returns the count used for plural translation selection.
func (e SimpleEvent) PluralCount() (int, bool) { return e.pluralCount, e.plural }

// With returns a copy of e with attrs appended.
func (e SimpleEvent) With(attrs ...Attr) SimpleEvent {
	e.attrs = append(cloneAttrs(e.attrs), attrs...)
	return e
}

// WithPluralCount returns a copy of e marked as pluralized with count.
func (e SimpleEvent) WithPluralCount(count int) SimpleEvent {
	e.pluralCount = count
	e.plural = true
	return e
}

func cloneAttrs(attrs []Attr) []Attr {
	if len(attrs) == 0 {
		return nil
	}
	cp := make([]Attr, len(attrs))
	copy(cp, attrs)
	return cp
}
