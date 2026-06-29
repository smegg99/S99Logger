// attr.go

package s99logger

import "time"

// Attr is a single structured key/value pair attached to an event. Attrs are
// emitted in the log record and are also used as template data by translators.
type Attr struct {
	Key   string
	Value any
}

// String returns a string attribute.
func String(key, value string) Attr { return Attr{Key: key, Value: value} }

// Int returns an integer attribute.
func Int(key string, value int) Attr { return Attr{Key: key, Value: value} }

// Bool returns a boolean attribute.
func Bool(key string, value bool) Attr { return Attr{Key: key, Value: value} }

// Duration returns a duration attribute. The duration renders as text (e.g.
// "1.5s") both in log output and in translation templates.
func Duration(key string, value time.Duration) Attr { return Attr{Key: key, Value: value} }

// Err returns an attribute carrying an error message under the "error" key.
func Err(err error) Attr {
	if err == nil {
		return Attr{Key: "error", Value: ""}
	}
	return Attr{Key: "error", Value: err.Error()}
}

// Any returns an attribute with an arbitrary value.
func Any(key string, value any) Attr { return Attr{Key: key, Value: value} }
