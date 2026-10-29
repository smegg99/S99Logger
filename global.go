// global.go

package s99logger

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
)

// The process-wide Logger, swappable at runtime with SetDefault.
var (
	mu      sync.RWMutex
	global  = New(NewConsoleSink(os.Stderr), Options{MinLevel: LevelDebug})
	closers []io.Closer
)

// Default returns the process-wide Logger; before SetDefault it writes to stderr at debug level.
func Default() *Logger {
	mu.RLock()
	defer mu.RUnlock()
	return global
}

// SetDefault installs l and its sinkClosers, closing the previous closers; a nil l is ignored.
func SetDefault(l *Logger, sinkClosers ...io.Closer) {
	if l == nil {
		return
	}
	mu.Lock()
	old := closers
	global = l
	closers = sinkClosers
	mu.Unlock()

	for _, c := range old {
		_ = c.Close()
	}
}

// Close releases the default Logger's closers; safe to call more than once.
func Close() error {
	mu.Lock()
	cs := closers
	closers = nil
	mu.Unlock()

	var firstErr error
	for _, c := range cs {
		if err := c.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// With is Default().With.
func With(attrs ...Attr) *Logger { return Default().With(attrs...) }

// Debug logs event at debug level on the default Logger.
func Debug(event Event, ctx ...context.Context) { Default().Debug(event, ctx...) }

// Info logs event at info level on the default Logger.
func Info(event Event, ctx ...context.Context) { Default().Info(event, ctx...) }

// Warn logs event at warn level on the default Logger.
func Warn(event Event, ctx ...context.Context) { Default().Warn(event, ctx...) }

// Error logs event at error level on the default Logger.
func Error(event Event, ctx ...context.Context) { Default().Error(event, ctx...) }

// Fatal logs at error level, closes the sinks and exits with status 1 without running deferred functions.
func Fatal(event Event, ctx ...context.Context) {
	Default().Error(event, ctx...)
	_ = Close()
	os.Exit(1)
}

// ParseLevel reads debug, info, warn, warning or error in any case; ok is false otherwise.
func ParseLevel(name string) (level Level, ok bool) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "debug":
		return LevelDebug, true
	case "info":
		return LevelInfo, true
	case "warn", "warning":
		return LevelWarn, true
	case "error":
		return LevelError, true
	}
	return LevelInfo, false
}
