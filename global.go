// global.go

package s99logger

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
)

// The package keeps a single process-wide Logger so applications can log from
// anywhere without threading a *Logger through every call. It is safe for
// concurrent use and can be swapped at runtime (e.g. once configuration is
// loaded) via SetDefault.
var (
	mu      sync.RWMutex
	global  = New(NewConsoleSink(os.Stderr), Options{MinLevel: LevelDebug})
	closers []io.Closer
)

// Default returns the process-wide Logger. It is never nil; before SetDefault
// is called it writes to stderr at debug level.
func Default() *Logger {
	mu.RLock()
	defer mu.RUnlock()
	return global
}

// SetDefault installs l as the process-wide Logger. Any sinks that need cleanup
// at shutdown (such as file sinks) should be passed as sinkClosers; they are
// closed by Close. Closers registered by a previous SetDefault are closed here,
// so reconfiguring the logger releases the old file handles. A nil l is ignored.
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

// Close releases the closers registered with the current default Logger. Call
// it during shutdown. It is safe to call more than once.
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

// With returns a child of the default Logger that attaches attrs to every
// record it emits. It is a shorthand for Default().With and is handy for
// request- or scope-specific loggers. See Logger.With.
func With(attrs ...Attr) *Logger { return Default().With(attrs...) }

// Debug logs event at debug level on the default Logger.
func Debug(event Event, ctx ...context.Context) { Default().Debug(event, ctx...) }

// Info logs event at info level on the default Logger.
func Info(event Event, ctx ...context.Context) { Default().Info(event, ctx...) }

// Warn logs event at warn level on the default Logger.
func Warn(event Event, ctx ...context.Context) { Default().Warn(event, ctx...) }

// Error logs event at error level on the default Logger.
func Error(event Event, ctx ...context.Context) { Default().Error(event, ctx...) }

// Fatal logs event at error level on the default Logger, releases its sinks so
// the record reaches disk, and terminates the process with status 1. Deferred
// functions do not run; reserve it for startup failures the application cannot
// continue past.
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
