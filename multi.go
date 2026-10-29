// multi.go

package s99logger

import (
	"context"
	"errors"
)

// MultiSink writes each record to every sink and joins their errors.
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
