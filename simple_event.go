// simple_event.go

package s99logger

// SimpleEvent is an immutable Event for events that need no struct of their own.
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

// NewPluralEvent returns an event whose count selects the plural form.
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

func cloneAttrs(attrs []Attr) []Attr {
	if len(attrs) == 0 {
		return nil
	}
	cp := make([]Attr, len(attrs))
	copy(cp, attrs)
	return cp
}
