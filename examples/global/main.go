// examples/global/main.go
package main

import (
	"errors"

	"github.com/smegg99/s99logger"
	"github.com/smegg99/s99logger/rotation"
)

func serverStarting(port int) s99logger.Event {
	return s99logger.NewEvent("server_starting", s99logger.Int("port", port))
}

func main() {
	// rotation.Configure builds a console logger (plus a rotating JSON file sink
	// when EnableFiles is set) and installs it as the process-wide default. Set
	// EnableFiles + Directory to also write rotating files. Until Configure runs,
	// the package-level functions already work, writing to stderr at debug level.
	if err := rotation.Configure(rotation.Config{
		Service: "api",
		Level:   "info",
	}); err != nil {
		panic(err)
	}
	defer s99logger.Close()

	// Log from anywhere via the package-level functions, no *Logger to thread
	// through your call stack.
	s99logger.Info(serverStarting(8080))
	s99logger.Warn(s99logger.NewEvent("cache_miss", s99logger.String("key", "user:42")))
	s99logger.Error(s99logger.NewEvent("db_unreachable", s99logger.Err(errors.New("connection refused"))))

	// Default() exposes the installed logger when you need a *Logger value, e.g.
	// to derive a request-scoped child with With.
	reqLog := s99logger.Default().With(s99logger.String("request_id", "abc123"))
	reqLog.Info(s99logger.NewEvent("request_handled"))
}
