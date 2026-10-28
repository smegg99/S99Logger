// slog.go

package s99logger

import (
	"context"
	"log/slog"
	"slices"
)

// Enabled reports whether a record at level reaches the sink.
func (l *Logger) Enabled(level Level) bool {
	return l != nil && !isNilInterface(l.sink) && level >= l.minLevel
}

// SlogHandler returns a slog.Handler that logs each record as event id, the slog message under "detail".
func (l *Logger) SlogHandler(id MessageID) slog.Handler {
	return &slogHandler{logger: l, id: id}
}

// slogHandler adapts slog records to events; groups become dotted key prefixes.
type slogHandler struct {
	logger *Logger
	id     MessageID
	attrs  []Attr
	group  string
}

func (h *slogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.logger.Enabled(levelFromSlog(level))
}

func (h *slogHandler) Handle(ctx context.Context, r slog.Record) error {
	level := levelFromSlog(r.Level)
	if !h.logger.Enabled(level) {
		return nil
	}
	attrs := append(slices.Clone(h.attrs), String("detail", r.Message))
	r.Attrs(func(a slog.Attr) bool {
		attrs = appendSlogAttr(attrs, h.group, a)
		return true
	})
	h.logger.log(level, NewEvent(h.id, attrs...), ctx)
	return nil
}

func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := *h
	next.attrs = slices.Clone(h.attrs)
	for _, a := range attrs {
		next.attrs = appendSlogAttr(next.attrs, h.group, a)
	}
	return &next
}

func (h *slogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	next := *h
	next.attrs = slices.Clone(h.attrs)
	next.group = h.group + name + "."
	return &next
}

// levelFromSlog maps a slog level to the nearest level at or below it.
func levelFromSlog(level slog.Level) Level {
	switch {
	case level >= slog.LevelError:
		return LevelError
	case level >= slog.LevelWarn:
		return LevelWarn
	case level >= slog.LevelInfo:
		return LevelInfo
	default:
		return LevelDebug
	}
}

// appendSlogAttr flattens a, writing an error as its text.
func appendSlogAttr(attrs []Attr, prefix string, a slog.Attr) []Attr {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return attrs
	}
	if a.Value.Kind() == slog.KindGroup {
		if a.Key != "" {
			prefix += a.Key + "."
		}
		for _, child := range a.Value.Group() {
			attrs = appendSlogAttr(attrs, prefix, child)
		}
		return attrs
	}
	value := a.Value.Any()
	if err, ok := value.(error); ok {
		value = err.Error()
	}
	return append(attrs, Any(prefix+a.Key, value))
}
