// tests/helpers_test.go
package tests

import (
	"context"

	"github.com/smegg99/s99logger"
)

// event is a configurable Event used across the test suite. Set plural to true
// to exercise pluralized translation.
type event struct {
	id     s99logger.MessageID
	attrs  []s99logger.Attr
	count  int
	plural bool
}

func (e event) MessageID() s99logger.MessageID { return e.id }
func (e event) Attrs() []s99logger.Attr        { return e.attrs }
func (e event) PluralCount() (int, bool)       { return e.count, e.plural }

// captureSink records every record it receives.
type captureSink struct {
	records []s99logger.Record
}

func (s *captureSink) Write(_ context.Context, rec s99logger.Record) error {
	s.records = append(s.records, rec)
	return nil
}

// stubTranslator returns a fixed message, or an error when err is set.
type stubTranslator struct {
	msg string
	err error
}

func (t stubTranslator) Translate(_ context.Context, _ string, _ s99logger.Event) (string, error) {
	return t.msg, t.err
}
