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
