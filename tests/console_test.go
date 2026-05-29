// tests/console_test.go
package tests

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/smegg99/s99logger"
)

func TestConsoleSinkReadableLine(t *testing.T) {
	var buf bytes.Buffer
	sink := s99logger.NewConsoleSink(&buf).WithColor(false)
	logger := s99logger.New(sink, s99logger.Options{
		Service:    "demo",
		Translator: stubTranslator{msg: "hello world"},
		Clock:      func() time.Time { return time.Unix(0, 0) },
	})

	logger.Info(event{
		id:    "greeting",
		attrs: []s99logger.Attr{s99logger.String("user", "alice")},
	})

	line := buf.String()
	for _, want := range []string{"INFO", "[demo]", "hello world", "user=alice", "id=greeting"} {
		if !strings.Contains(line, want) {
			t.Errorf("console line missing %q\ngot: %s", want, line)
		}
	}
	if strings.Contains(line, "service=demo") {
		t.Errorf("service should render as a prefix, got: %s", line)
	}
	if serviceAt, levelAt := strings.Index(line, "[demo]"), strings.Index(line, "INFO"); serviceAt < 0 || levelAt < 0 || serviceAt > levelAt {
		t.Errorf("service prefix should appear before level prefix, got: %s", line)
	}
}

func TestConsoleSinkColorWhenEnabled(t *testing.T) {
	var buf bytes.Buffer
	sink := s99logger.NewConsoleSink(&buf).WithColor(true)
	logger := s99logger.New(sink, s99logger.Options{})

	logger.Error(event{id: "boom"})

	if !strings.Contains(buf.String(), "\x1b[") {
		t.Errorf("expected ANSI color codes, got: %q", buf.String())
	}
}

func TestConsoleSinkCustomLevelPrefixAndColor(t *testing.T) {
	var buf bytes.Buffer
	sink := s99logger.NewConsoleSink(&buf).
		WithColor(true).
		WithLevelPrefix(s99logger.LevelInfo, "OK").
		WithLevelColor(s99logger.LevelInfo, s99logger.ANSIBrightMagenta).
		WithFieldColor("")
	logger := s99logger.New(sink, s99logger.Options{
		Clock: func() time.Time { return time.Unix(0, 0) },
	})

	logger.Info(event{id: "ready", attrs: []s99logger.Attr{s99logger.String("mode", "custom")}})

	line := buf.String()
	for _, want := range []string{s99logger.ANSIBrightMagenta, "OK", "mode=custom"} {
		if !strings.Contains(line, want) {
			t.Errorf("console line missing %q\ngot: %q", want, line)
		}
	}
	if strings.Contains(line, "INFO") {
		t.Errorf("default prefix leaked into customized line: %q", line)
	}
}

func TestConsoleSinkCustomServicePrefix(t *testing.T) {
	var buf bytes.Buffer
	sink := s99logger.NewConsoleSink(&buf).
		WithColor(true).
		WithServicePrefixStyle("<", ">").
		WithServiceColor(s99logger.ANSIBrightWhite)
	logger := s99logger.New(sink, s99logger.Options{
		Service: "billing",
		Clock:   func() time.Time { return time.Unix(0, 0) },
	})

	logger.Info(event{id: "paid"})

	line := buf.String()
	for _, want := range []string{s99logger.ANSIBrightWhite, "<billing>", "paid"} {
		if !strings.Contains(line, want) {
			t.Errorf("console line missing %q\ngot: %q", want, line)
		}
	}
	if strings.Contains(line, "service=billing") {
		t.Errorf("service should not render as a field, got: %q", line)
	}
}

func TestConsoleSinkLocalizedLevelPrefixes(t *testing.T) {
	var buf bytes.Buffer
	sink := s99logger.NewConsoleSink(&buf).
		WithColor(true).
		WithLocalizedLevelPrefix("pl", s99logger.LevelWarn, "UWAGA").
		WithLocalizedLevelColor("pl", s99logger.LevelWarn, s99logger.ANSIBrightBlue)
	logger := s99logger.New(sink, s99logger.Options{
		Clock: func() time.Time { return time.Unix(0, 0) },
	})

	ctx := s99logger.ContextWithLanguage(context.Background(), "pl-PL")
	logger.Warn(event{id: "slow"}, ctx)

	line := buf.String()
	if !strings.Contains(line, "UWAGA") {
		t.Errorf("localized prefix missing\ngot: %q", line)
	}
	if !strings.Contains(line, s99logger.ANSIBrightBlue) {
		t.Errorf("localized color missing\ngot: %q", line)
	}
	if strings.Contains(line, "WARN") {
		t.Errorf("default prefix leaked into localized line: %q", line)
	}
}

func TestMultiSinkWritesToAll(t *testing.T) {
	var console, jsonOut bytes.Buffer
	sink := s99logger.MultiSink(
		s99logger.NewConsoleSink(&console).WithColor(false),
		s99logger.NewJSONSink(&jsonOut),
	)
	logger := s99logger.New(sink, s99logger.Options{})

	logger.Info(event{id: "evt"})

	if !strings.Contains(console.String(), "evt") {
		t.Errorf("console sink got no output: %q", console.String())
	}
	if !strings.Contains(jsonOut.String(), `"message_id":"evt"`) {
		t.Errorf("json sink got no output: %q", jsonOut.String())
	}
}

func TestConsoleSinkConcurrentConfigurationAndWrites(t *testing.T) {
	var buf bytes.Buffer
	sink := s99logger.NewConsoleSink(&buf).WithColor(false)
	logger := s99logger.New(sink, s99logger.Options{
		Clock: func() time.Time { return time.Unix(0, 0) },
	})

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			logger.Info(event{id: "evt"})
		}()
		go func() {
			defer wg.Done()
			sink.WithLevelPrefix(s99logger.LevelInfo, "OK")
			sink.WithLevelColor(s99logger.LevelInfo, s99logger.ANSIGreen)
		}()
	}
	wg.Wait()
}
