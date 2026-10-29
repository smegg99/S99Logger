// rotation/configure.go

package rotation

import (
	"fmt"
	"io"
	"os"

	"github.com/smegg99/s99logger"
)

// Config is Configure's input: a console logger with an optional rotating JSON file.
type Config struct {
	// Service is the service name attached to every record. Defaults to "app".
	Service string
	// Level is the minimum level name (see s99logger.ParseLevel), info when empty. Ignored when
	// Verbose is set.
	Level string
	// Verbose forces debug level regardless of Level.
	Verbose bool
	// NoColor disables ANSI colors on the console sink.
	NoColor bool

	// EnableFiles turns on the rotating file sink. When false (or Directory is
	// empty) only the console sink is used.
	EnableFiles bool
	// Directory, Filename, MaxSizeMB, MaxBackups, MaxAgeDays, Compression and
	// LocalTime configure the file sink; see Options for details.
	Directory   string
	Filename    string
	MaxSizeMB   int
	MaxBackups  int
	MaxAgeDays  int
	Compression string
	LocalTime   bool
}

// Configure builds a Logger from cfg and installs it as the default, registering the file sink for s99logger.Close.
func Configure(cfg Config) error {
	level := s99logger.LevelInfo
	switch {
	case cfg.Verbose:
		level = s99logger.LevelDebug
	case cfg.Level != "":
		var ok bool
		if level, ok = s99logger.ParseLevel(cfg.Level); !ok {
			return fmt.Errorf("%w: %q", errLevelUnknown, cfg.Level)
		}
	}

	console := s99logger.NewConsoleSink(os.Stderr)
	if cfg.NoColor {
		console.WithColor(false)
	}
	sinks := []s99logger.Sink{console}

	var closers []io.Closer
	if cfg.EnableFiles && cfg.Directory != "" {
		file, err := New(Options{
			Directory:   cfg.Directory,
			Filename:    cfg.Filename,
			MaxSizeMB:   cfg.MaxSizeMB,
			MaxBackups:  cfg.MaxBackups,
			MaxAgeDays:  cfg.MaxAgeDays,
			Compression: cfg.Compression,
			LocalTime:   cfg.LocalTime,
		})
		if err != nil {
			return err
		}
		sinks = append(sinks, file)
		closers = append(closers, file)
	}

	service := cfg.Service
	if service == "" {
		service = "app"
	}

	logger := s99logger.New(s99logger.MultiSink(sinks...), s99logger.Options{
		Service:  service,
		MinLevel: level,
	})
	s99logger.SetDefault(logger, closers...)
	return nil
}
