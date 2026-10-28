// filter.go

package s99logger

import "context"

// MinLevel returns a Sink that passes on records at min or above, so one logger can feed
// sinks at different levels. Run the logger at the lowest level any sink needs.
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
