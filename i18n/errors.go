package i18n

import "errors"

// Internal sentinel errors reused across the package.
var (
	errFilesystemRequired  = errors.New("FS is required")
	errMissingTranslations = errors.New("missing translations")
)
