// logger.go

package s99logger

import "time"

// Logger emits typed events to a Sink, optionally translated; safe for concurrent use when its Sink is.
type Logger struct {
	sink       Sink
	translator Translator
	lang       string
	service    string
	minLevel   Level
	base       []Attr
	clock      func() time.Time
}

// New returns a Logger writing to sink; a nil sink discards records.
func New(sink Sink, opts Options) *Logger {
	clock := opts.Clock
	if clock == nil {
		clock = time.Now
	}
	return &Logger{
		sink:       sink,
		translator: opts.Translator,
		lang:       opts.Language,
		service:    opts.Service,
		minLevel:   opts.MinLevel,
		clock:      clock,
	}
}

// With returns a child Logger that adds attrs to every record; the receiver is unchanged.
func (l *Logger) With(attrs ...Attr) *Logger {
	if len(attrs) == 0 {
		return l
	}
	cp := *l
	cp.base = make([]Attr, 0, len(l.base)+len(attrs))
	cp.base = append(cp.base, l.base...)
	cp.base = append(cp.base, attrs...)
	return &cp
}

// mergeAttrs prepends the base attrs from With to the event's attrs.
func (l *Logger) mergeAttrs(eventAttrs []Attr) []Attr {
	if len(l.base) == 0 {
		return eventAttrs
	}
	merged := make([]Attr, 0, len(l.base)+len(eventAttrs))
	merged = append(merged, l.base...)
	merged = append(merged, eventAttrs...)
	return merged
}
