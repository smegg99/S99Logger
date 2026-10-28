// tests/slog_test.go
package tests

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/smegg99/s99logger"
)

// Enabled answers from the minimum level and is false without a sink.
func TestEnabled(t *testing.T) {
	l := s99logger.New(&captureSink{}, s99logger.Options{MinLevel: s99logger.LevelWarn})
	if l.Enabled(s99logger.LevelInfo) || !l.Enabled(s99logger.LevelError) {
		t.Fatal("Enabled does not follow MinLevel")
	}
	if s99logger.New(nil, s99logger.Options{}).Enabled(s99logger.LevelError) {
		t.Fatal("a logger without a sink reports enabled")
	}
	var none *s99logger.Logger
	if none.Enabled(s99logger.LevelError) {
		t.Fatal("a nil logger reports enabled")
	}
}

// A slog record becomes one event with the message as detail and groups as dotted keys.
func TestSlogHandler(t *testing.T) {
	sink := &captureSink{}
	l := s99logger.New(sink, s99logger.Options{MinLevel: s99logger.LevelInfo})
	log := slog.New(l.SlogHandler("runtime")).With("component", "tray").WithGroup("req")

	log.Debug("dropped")
	log.Warn("slow", "ms", 250, slog.Group("peer", "host", "a"), "err", errors.New("boom"))

	if len(sink.records) != 1 {
		t.Fatalf("records = %d, want 1", len(sink.records))
	}
	rec := sink.records[0]
	if rec.Level != s99logger.LevelWarn || rec.MessageID != "runtime" {
		t.Fatalf("record = %+v", rec)
	}
	got := map[string]any{}
	for _, a := range rec.Attrs {
		got[a.Key] = a.Value
	}
	want := map[string]any{"component": "tray", "detail": "slow", "req.ms": int64(250), "req.peer.host": "a", "req.err": "boom"}
	for key, value := range want {
		if got[key] != value {
			t.Fatalf("attr %s = %#v, want %#v (all: %#v)", key, got[key], value, got)
		}
	}
	if !l.SlogHandler("runtime").Enabled(context.Background(), slog.LevelWarn) {
		t.Fatal("handler disabled for warn")
	}
}
