// json.go

package s99logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

// reserved keys are written from the record itself and cannot be overridden by
// event attributes.
var reservedKeys = map[string]struct{}{
	"time":       {},
	"level":      {},
	"service":    {},
	"lang":       {},
	"message":    {},
	"message_id": {},
}

// JSONSink writes records as newline-delimited JSON to an io.Writer. It is
// safe for concurrent use.
type JSONSink struct {
	mu  sync.Mutex
	out io.Writer
}

// NewJSONSink returns a JSONSink writing to w.
func NewJSONSink(w io.Writer) *JSONSink {
	if w == nil {
		w = io.Discard
	}
	return &JSONSink{out: w}
}

// Write encodes rec as a single JSON object followed by a newline.
func (s *JSONSink) Write(_ context.Context, rec Record) error {
	obj := map[string]any{
		"time":       rec.Time.UTC().Format(time.RFC3339Nano),
		"level":      rec.Level.String(),
		"message":    rec.Message,
		"message_id": string(rec.MessageID),
	}
	if rec.Service != "" {
		obj["service"] = rec.Service
	}
	if rec.Lang != "" {
		obj["lang"] = rec.Lang
	}
	for _, a := range rec.Attrs {
		if _, ok := reservedKeys[a.Key]; ok {
			continue
		}
		obj[a.Key] = jsonValue(a.Value)
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("%w: %w", errMarshalRecord, err)
	}
	data = append(data, '\n')

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := s.out.Write(data); err != nil {
		return fmt.Errorf("%w: %w", errWriteJSONRecord, err)
	}
	return nil
}

// jsonValue normalizes attribute values for JSON output. Durations are written
// as text so they match how they render in translation templates.
func jsonValue(v any) any {
	if d, ok := v.(time.Duration); ok {
		return d.String()
	}
	return v
}
