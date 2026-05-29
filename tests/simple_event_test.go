package tests

import (
	"testing"

	"github.com/smegg99/s99logger"
)

func TestNewEvent(t *testing.T) {
	event := s99logger.NewEvent("server_starting", s99logger.Int("port", 8080))

	if event.MessageID() != "server_starting" {
		t.Fatalf("MessageID = %q, want server_starting", event.MessageID())
	}
	if count, ok := event.PluralCount(); ok || count != 0 {
		t.Fatalf("PluralCount = (%d, %v), want (0, false)", count, ok)
	}

	attrs := event.Attrs()
	if len(attrs) != 1 || attrs[0].Key != "port" || attrs[0].Value != 8080 {
		t.Fatalf("Attrs = %#v, want port=8080", attrs)
	}
}

func TestNewPluralEvent(t *testing.T) {
	event := s99logger.NewPluralEvent("jobs_processed", 3, s99logger.Int("count", 3))

	count, ok := event.PluralCount()
	if !ok || count != 3 {
		t.Fatalf("PluralCount = (%d, %v), want (3, true)", count, ok)
	}
}

func TestSimpleEventAttrsAreCopied(t *testing.T) {
	attrs := []s99logger.Attr{s99logger.String("key", "before")}
	event := s99logger.NewEvent("evt", attrs...)
	attrs[0] = s99logger.String("key", "after")

	got := event.Attrs()
	if got[0].Value != "before" {
		t.Fatalf("constructor retained caller attrs: %#v", got)
	}

	got[0] = s99logger.String("key", "mutated")
	again := event.Attrs()
	if again[0].Value != "before" {
		t.Fatalf("Attrs exposed internal attrs: %#v", again)
	}
}
