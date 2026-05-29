package rotation

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultFilename = "app.log"
	defaultRotateAt = "00:00"
)

// Options configures a JSON log file that is rotated once per day. Directory
// is required so applications explicitly choose where logs live.
type Options struct {
	// Directory is created if needed and receives the active and rotated logs.
	Directory string
	// Filename is the active log file name inside Directory. It defaults to
	// "app.log" and must not be an absolute path.
	Filename string
	// RotateAt is the HH:MM time used for daily rotation. It defaults to
	// midnight. timberjack uses UTC unless LocalTime is true.
	RotateAt string
	// MaxSizeMB is timberjack's size rotation limit in megabytes. The
	// timberjack default is used when this is zero.
	MaxSizeMB int
	// MaxAgeDays removes rotated logs older than this many days. Zero keeps
	// files based only on MaxBackups.
	MaxAgeDays int
	// MaxBackups is the maximum number of rotated files to retain. Zero keeps
	// all backups unless MaxAgeDays removes them.
	MaxBackups int
	// Compression may be "none", "gzip", or "zstd". Empty uses timberjack's
	// default of no compression.
	Compression string
	// LocalTime makes rotation times and backup filenames use local time instead
	// of UTC.
	LocalTime bool
	// FileMode configures permissions for newly created files. Zero uses
	// timberjack's default.
	FileMode os.FileMode
}

func normalizeOptions(opts Options) (Options, error) {
	if opts.Directory == "" {
		return Options{}, errDirectoryRequired
	}
	if opts.MaxSizeMB < 0 {
		return Options{}, errMaxSizeNegative
	}
	if opts.MaxAgeDays < 0 {
		return Options{}, errMaxAgeNegative
	}
	if opts.MaxBackups < 0 {
		return Options{}, errMaxBackupsNegative
	}

	if opts.Filename == "" {
		opts.Filename = defaultFilename
	}
	if filepath.IsAbs(opts.Filename) {
		return Options{}, errFilenameAbsolute
	}
	opts.Filename = filepath.Clean(opts.Filename)
	if opts.Filename == "." || opts.Filename == ".." || opts.Filename != filepath.Base(opts.Filename) {
		return Options{}, errFilenamePath
	}

	if opts.RotateAt == "" {
		opts.RotateAt = defaultRotateAt
	}
	if _, err := time.Parse("15:04", opts.RotateAt); err != nil {
		return Options{}, fmt.Errorf("%w: %v", errRotateAtInvalid, err)
	}

	return opts, nil
}
