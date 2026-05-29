// tests/logger_test.go
package tests

import (
	"errors"
	"testing"

	"github.com/smegg99/s99logger"
)

func TestFallbackToMessageIDWithoutTranslator(t *testing.T) {
	sink := &captureSink{}
	logger := s99logger.New(sink, s99logger.Options{})

	logger.Info(event{id: "app_started"})

	rec := sink.records[0]
	if rec.Message != "app_started" || rec.MessageID != "app_started" {
		t.Errorf("got message=%q id=%q, want both %q", rec.Message, rec.MessageID, "app_started")
	}
}

func TestLocalizedMessageWithTranslator(t *testing.T) {
	sink := &captureSink{}
	logger := s99logger.New(sink, s99logger.Options{
		Translator: stubTranslator{msg: "Uruchomiono"},
	})

	logger.Info(event{id: "app_started"})

	rec := sink.records[0]
	if rec.Message != "Uruchomiono" {
		t.Errorf("message = %q, want localized %q", rec.Message, "Uruchomiono")
	}
	if rec.MessageID != "app_started" {
		t.Errorf("message_id = %q, want stable %q", rec.MessageID, "app_started")
	}
}

func TestTranslationFailureStillLogs(t *testing.T) {
	sink := &captureSink{}
	logger := s99logger.New(sink, s99logger.Options{
		Translator: stubTranslator{err: errors.New("missing")},
	})

	logger.Error(event{id: "login_failed"})

	if len(sink.records) != 1 {
		t.Fatalf("expected log to still be emitted, got %d records", len(sink.records))
	}
	if sink.records[0].Message != "login_failed" {
		t.Errorf("message = %q, want fallback to id", sink.records[0].Message)
	}
}

func TestLevelsMapToMethods(t *testing.T) {
	sink := &captureSink{}
	logger := s99logger.New(sink, s99logger.Options{})
	logger.Debug(event{id: "m"})
	logger.Info(event{id: "m"})
	logger.Warn(event{id: "m"})
	logger.Error(event{id: "m"})

	want := []s99logger.Level{s99logger.LevelDebug, s99logger.LevelInfo, s99logger.LevelWarn, s99logger.LevelError}
	for i, lvl := range want {
		if sink.records[i].Level != lvl {
			t.Errorf("record %d level = %v, want %v", i, sink.records[i].Level, lvl)
		}
	}
}

func TestNilSinkDiscards(t *testing.T) {
	logger := s99logger.New(nil, s99logger.Options{})

	logger.Info(s99logger.NewEvent("discarded"))
}

func TestNilEventIsIgnored(t *testing.T) {
	sink := &captureSink{}
	logger := s99logger.New(sink, s99logger.Options{})

	logger.Info(nil)

	if len(sink.records) != 0 {
		t.Fatalf("nil event produced %d records, want 0", len(sink.records))
	}
}

func TestTypedNilEventIsIgnored(t *testing.T) {
	sink := &captureSink{}
	logger := s99logger.New(sink, s99logger.Options{})
	var event *typedNilEvent

	logger.Info(event)

	if len(sink.records) != 0 {
		t.Fatalf("typed nil event produced %d records, want 0", len(sink.records))
	}
}

type typedNilEvent struct{}

func (*typedNilEvent) MessageID() s99logger.MessageID { return "nil" }
func (*typedNilEvent) Attrs() []s99logger.Attr        { return nil }
func (*typedNilEvent) PluralCount() (int, bool)       { return 0, false }
