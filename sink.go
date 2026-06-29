// sink.go

package s99logger

import (
	"context"
	"time"
)

// Record is a fully resolved log entry handed to a Sink. Message holds the
// localized text (or the message id, when translation is unavailable), while
// MessageID always holds the stable identifier.
type Record struct {
	Time      time.Time
	Level     Level
	Service   string
	Lang      string
	MessageID MessageID
	Message   string
	Attrs     []Attr
}

// Sink writes resolved log records to a destination. Implementations must be
// safe for concurrent use. The returned error is for the sink's own use (e.g.
// custom sinks and tests); the Logger discards it, like slog.
type Sink interface {
	Write(ctx context.Context, rec Record) error
}
