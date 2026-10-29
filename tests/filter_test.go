// tests/filter_test.go

package tests

import (
	"testing"

	"github.com/smegg99/s99logger"
)

func TestMinLevelFeedsSinksAtDifferentLevels(t *testing.T) {
	all, warn := &captureSink{}, &captureSink{}
	log := s99logger.New(s99logger.MultiSink(all, s99logger.MinLevel(warn, s99logger.LevelWarn)), s99logger.Options{MinLevel: s99logger.LevelDebug})

	log.Debug(event{id: "debug"})
	log.Info(event{id: "info"})
	log.Warn(event{id: "warn"})
	log.Error(event{id: "error"})

	if len(all.records) != 4 {
		t.Fatalf("unfiltered sink got %d records, want 4", len(all.records))
	}
	if len(warn.records) != 2 || warn.records[0].MessageID != "warn" || warn.records[1].MessageID != "error" {
		t.Fatalf("warn sink got %+v, want warn and error", warn.records)
	}
}

func TestMinLevelOfNilIsNil(t *testing.T) {
	if s99logger.MinLevel(nil, s99logger.LevelInfo) != nil {
		t.Fatal("a filter over no sink must be no sink")
	}
}
