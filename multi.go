package s99logger

import (
	"context"
	"errors"
)

// MultiSink returns a Sink that writes each record to every given sink. This is
// the usual way to send colored lines to the terminal and JSON to a file at the
// same time. Sink errors are joined; the Logger discards them anyway.
func MultiSink(sinks ...Sink) Sink {
	clean := make([]Sink, 0, len(sinks))
	for _, sink := range sinks {
		if !isNilInterface(sink) {
			clean = append(clean, sink)
		}
	}
	return multiSink{sinks: clean}
}

type multiSink struct {
	sinks []Sink
}

// Write writes rec to every underlying sink, returning the joined errors.
func (m multiSink) Write(ctx context.Context, rec Record) error {
	var errs []error
	for _, s := range m.sinks {
		if isNilInterface(s) {
			continue
		}
		if err := s.Write(ctx, rec); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
