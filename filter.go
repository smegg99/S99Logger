// filter.go

package s99logger

import "context"

// MinLevel passes on records at min or above, so one logger at the lowest level can feed sinks at different levels.
func MinLevel(sink Sink, min Level) Sink {
	if isNilInterface(sink) {
		return nil
	}
	return levelSink{sink: sink, min: min}
}

type levelSink struct {
	sink Sink
	min  Level
}

func (s levelSink) Write(ctx context.Context, rec Record) error {
	if rec.Level < s.min {
		return nil
	}
	return s.sink.Write(ctx, rec)
}
