// tests/json_test.go
package tests

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/smegg99/s99logger"
)

func TestJSONSinkValidNDJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := s99logger.New(s99logger.NewJSONSink(&buf), s99logger.Options{
		Service: "svc",
		Clock:   func() time.Time { return time.Unix(0, 0).UTC() },
	})

	logger.Info(event{id: "first", attrs: []s99logger.Attr{s99logger.String("k", "v")}})
	logger.Error(event{id: "second"})

	lines := bytes.Split(bytes.TrimRight(buf.Bytes(), "\n"), []byte("\n"))
	if len(lines) != 2 {
		t.Fatalf("expected 2 NDJSON lines, got %d", len(lines))
	}

	var first map[string]any
	if err := json.Unmarshal(lines[0], &first); err != nil {
		t.Fatalf("line 1 invalid JSON: %v", err)
	}
	if first["message_id"] != "first" || first["service"] != "svc" || first["k"] != "v" || first["level"] != "info" {
		t.Errorf("unexpected line 1: %v", first)
	}

	var second map[string]any
	if err := json.Unmarshal(lines[1], &second); err != nil {
		t.Fatalf("line 2 invalid JSON: %v", err)
	}
	if second["level"] != "error" {
		t.Errorf("level = %v, want error", second["level"])
	}
}

func TestJSONSinkReservedKeysProtected(t *testing.T) {
	var buf bytes.Buffer
	logger := s99logger.New(s99logger.NewJSONSink(&buf), s99logger.Options{})

	// An attr named "message" must not clobber the resolved message.
	logger.Info(event{
		id:    "evt",
		attrs: []s99logger.Attr{s99logger.String("message", "hijack")},
	})

	var obj map[string]any
	if err := json.Unmarshal(bytes.TrimRight(buf.Bytes(), "\n"), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["message"] != "evt" {
		t.Errorf("message = %v, want evt (reserved key protected)", obj["message"])
	}
}

func TestJSONSinkNilWriterDiscards(t *testing.T) {
	logger := s99logger.New(s99logger.NewJSONSink(nil), s99logger.Options{})

	logger.Info(event{id: "discarded"})
}
