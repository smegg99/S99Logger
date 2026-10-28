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
		// One bad attribute (a NaN, a channel) must not lose the record: write it as text.
		for _, a := range rec.Attrs {
			if _, ok := reservedKeys[a.Key]; ok {
				continue
			}
			if _, err := json.Marshal(obj[a.Key]); err != nil {
				obj[a.Key] = fmt.Sprint(a.Value)
			}
		}
		if data, err = json.Marshal(obj); err != nil {
			return fmt.Errorf("%w: %w", errMarshalRecord, err)
		}
	}
	data = append(data, '\n')

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := s.out.Write(data); err != nil {
		return fmt.Errorf("%w: %w", errWriteJSONRecord, err)
	}
	return nil
}

// jsonValue writes durations and errors as text, as they render in translation templates.
func jsonValue(v any) any {
	switch v := v.(type) {
	case time.Duration:
		return v.String()
	case error:
		// Sprint, not Error: a typed nil error would panic.
		return fmt.Sprint(v)
	}
	return v
}
