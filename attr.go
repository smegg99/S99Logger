// attr.go

package s99logger

import "time"

// Attr is one key/value pair on an event, emitted in the record and passed to translators as template data.
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

// Duration returns a duration attribute, rendered as text such as "1.5s" in logs and templates.
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
