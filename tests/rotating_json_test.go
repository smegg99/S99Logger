// tests/rotating_json_test.go

package tests

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/smegg99/s99logger"
	"github.com/smegg99/s99logger/rotation"
)

func TestRotatingSinkWritesToConfiguredDirectory(t *testing.T) {
	logDir := filepath.Join(t.TempDir(), "logs")
	sink, err := rotation.New(rotation.Options{
		Directory:  logDir,
		Filename:   "service.log",
		LocalTime:  true,
		MaxBackups: 2,
	})
	if err != nil {
		t.Fatalf("rotation.New: %v", err)
	}
	t.Cleanup(func() {
		if err := sink.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	})

	logger := s99logger.New(sink, s99logger.Options{
		Service: "api",
		Clock: func() time.Time {
			return time.Date(2026, 5, 29, 10, 30, 0, 0, time.UTC)
		},
	})
	logger.Info(event{
		id:    "request_handled",
		attrs: []s99logger.Attr{s99logger.String("path", "/health")},
	})
	if err := sink.Sync(); err != nil {
		t.Fatalf("Sync: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(logDir, "service.log"))
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}

	var obj map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(data), &obj); err != nil {
		t.Fatalf("unmarshal log line %q: %v", data, err)
	}
	if obj["service"] != "api" {
		t.Errorf("service = %v, want api", obj["service"])
	}
	if obj["message_id"] != "request_handled" {
		t.Errorf("message_id = %v, want request_handled", obj["message_id"])
	}
	if obj["path"] != "/health" {
		t.Errorf("path = %v, want /health", obj["path"])
	}
}

func TestRotatingSinkRequiresDirectory(t *testing.T) {
	_, err := rotation.New(rotation.Options{})
	if err == nil {
		t.Fatal("expected an error when Directory is empty")
	}
}

func TestRotatingSinkRejectsPathFilename(t *testing.T) {
	_, err := rotation.New(rotation.Options{
		Directory: t.TempDir(),
		Filename:  filepath.Join("nested", "service.log"),
	})
	if err == nil {
		t.Fatal("expected an error when Filename includes a path")
	}
}

func TestRotatingSinkRejectsInvalidRotateAt(t *testing.T) {
	_, err := rotation.New(rotation.Options{
		Directory: t.TempDir(),
		RotateAt:  "24:00",
	})
	if err == nil {
		t.Fatal("expected an error when RotateAt is invalid")
	}
}

func TestRotationDropsWritesAfterClose(t *testing.T) {
	logDir := filepath.Join(t.TempDir(), "logs")
	sink, err := rotation.New(rotation.Options{Directory: logDir, Filename: "service.log"})
	if err != nil {
		t.Fatalf("rotation.New: %v", err)
	}
	logger := s99logger.New(sink, s99logger.Options{})
	logger.Info(event{id: "before"})
	if err := sink.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	path := filepath.Join(logDir, "service.log")
	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	logger.Info(event{id: "after"})

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if after.Size() != before.Size() {
		t.Fatalf("size %d -> %d, want a closed sink to drop the write", before.Size(), after.Size())
	}
}

func TestRotationCreatesItsDirectoryPrivate(t *testing.T) {
	logDir := filepath.Join(t.TempDir(), "logs")
	sink, err := rotation.New(rotation.Options{Directory: logDir, Filename: "service.log"})
	if err != nil {
		t.Fatalf("rotation.New: %v", err)
	}
	t.Cleanup(func() { _ = sink.Close() })

	info, err := os.Stat(logDir)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o750 {
		t.Fatalf("directory mode = %o, want 750", perm)
	}
}

func TestConfigureRejectsAnUnknownLevel(t *testing.T) {
	prev := s99logger.Default()
	t.Cleanup(func() { s99logger.SetDefault(prev) })

	if err := rotation.Configure(rotation.Config{Level: "verbose"}); err == nil {
		t.Fatal("Configure accepted the level \"verbose\"")
	}
	if s99logger.Default() != prev {
		t.Fatal("a refused Configure replaced the default logger")
	}
}
