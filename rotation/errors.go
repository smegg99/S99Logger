package rotation

import "errors"

// Internal sentinel errors reused across the package.
var (
	errDirectoryRequired  = errors.New("rotation: log directory is required")
	errMaxAgeNegative     = errors.New("rotation: MaxAgeDays must be >= 0")
	errMaxBackupsNegative = errors.New("rotation: MaxBackups must be >= 0")
	errMaxSizeNegative    = errors.New("rotation: MaxSizeMB must be >= 0")
	errFilenameAbsolute   = errors.New("rotation: Filename must be relative to Directory")
	errFilenamePath       = errors.New("rotation: Filename must be a file name, not a path")
	errRotateAtInvalid    = errors.New("rotation: RotateAt must use HH:MM 24-hour time")
)
