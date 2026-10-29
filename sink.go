// sink.go

package s99logger

import (
	"context"
	"time"
)

// Record is a resolved entry: Message is the translated text or the id, MessageID always the id.
type Record struct {
	Time      time.Time
	Level     Level
	Service   string
	Lang      string
	MessageID MessageID
	Message   string
	Attrs     []Attr
}

// Sink writes records somewhere; it must be safe for concurrent use, and the Logger discards its error.
type Sink interface {
	Write(ctx context.Context, rec Record) error
}
