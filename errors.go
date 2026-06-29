// errors.go

package s99logger

import "errors"

// Internal sentinel errors reused across the package.
var (
	errMarshalRecord      = errors.New("s99logger: marshal record")
	errWriteConsoleRecord = errors.New("s99logger: write console record")
	errWriteJSONRecord    = errors.New("s99logger: write JSON record")
)
