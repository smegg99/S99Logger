// rotation/rotation.go

package rotation

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/DeRuina/timberjack"

	"github.com/smegg99/s99logger"
)

// Sink writes newline-delimited JSON to a timberjack-managed file. Close it
// when the application shuts down so timberjack can release file handles and
// stop scheduled rotation work.
type Sink struct {
	json   *s99logger.JSONSink
	logger *timberjack.Logger

	// mu orders writes before Close: timberjack would reopen the file for a write after it.
	mu     sync.RWMutex
	closed bool
}

var _ s99logger.Sink = (*Sink)(nil)

// New returns a JSON sink backed by timberjack. The active log file is
// opts.Filename inside opts.Directory and is rotated every day.
func New(opts Options) (*Sink, error) {
	opts, err := normalizeOptions(opts)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(opts.Directory, 0o750); err != nil {
		return nil, fmt.Errorf("create log directory %q: %w", opts.Directory, err)
	}

	logger := &timberjack.Logger{
		Filename:    filepath.Join(opts.Directory, opts.Filename),
		MaxSize:     opts.MaxSizeMB,
		MaxAge:      opts.MaxAgeDays,
		MaxBackups:  opts.MaxBackups,
		LocalTime:   opts.LocalTime,
		Compression: opts.Compression,
		RotateAt:    []string{opts.RotateAt},
		FileMode:    opts.FileMode,
	}

	return &Sink{
		json:   s99logger.NewJSONSink(logger),
		logger: logger,
	}, nil
}

// Write encodes rec as a single JSON object and writes it to the active log
// file. After Close it drops rec.
func (s *Sink) Write(ctx context.Context, rec s99logger.Record) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return nil
	}
	return s.json.Write(ctx, rec)
}

// Close closes the underlying timberjack logger; later writes are dropped.
func (s *Sink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return s.logger.Close()
}

// Sync flushes the active log file when it has been opened.
func (s *Sink) Sync() error {
	return s.logger.Sync()
}

// Rotate forces an immediate rotation.
func (s *Sink) Rotate() error {
	return s.logger.Rotate()
}

// rotateWithReason forces an immediate rotation and tags the backup filename
// with reason after timberjack sanitizes it.
func (s *Sink) rotateWithReason(reason string) error {
	return s.logger.RotateWithReason(reason)
}
