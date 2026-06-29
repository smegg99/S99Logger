// rotation/configure.go

package rotation

import (
	"io"
	"os"

	"github.com/smegg99/s99logger"
)

// Config describes a console logger with an optional rotating JSON file sink.
// It is the high-level input to Configure, which most applications use instead
// of wiring sinks by hand.
type Config struct {
	// Service is the service name attached to every record. Defaults to "app".
	Service string
	// Level is the minimum level name (see s99logger.ParseLevel). Ignored when
	// Verbose is set.
	Level string
	// Verbose forces debug level regardless of Level.
	Verbose bool
	// NoColor disables ANSI colors on the console sink.
	NoColor bool

	// EnableFiles turns on the rotating file sink. When false (or Directory is
	// empty) only the console sink is used.
	EnableFiles bool
	// Directory, Filename, MaxSizeMB, MaxBackups, MaxAgeDays configure the file
	// sink; see Options for details.
	Directory  string
	Filename   string
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
}

// Configure builds a Logger from cfg and installs it as the s99logger default,
// registering the file sink (if any) so s99logger.Close releases it at
// shutdown. It returns any error from setting up the file sink.
func Configure(cfg Config) error {
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
			Compression: "gzip",
			LocalTime:   true,
		})
		if err != nil {
			return err
		}
		sinks = append(sinks, file)
		closers = append(closers, file)
	}

	level := s99logger.ParseLevel(cfg.Level)
	if cfg.Verbose {
		level = s99logger.LevelDebug
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
