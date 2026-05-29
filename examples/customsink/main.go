// examples/customsink/main.go
//
// Implementing your own Sink. A Sink just receives resolved Records, so you
// can send logs anywhere: a counter, a metrics system, a database, etc.

package main

import (
	"context"
	"fmt"

	"github.com/smegg99/s99logger"
)

// countingSink prints each record and tallies how many were seen per level.
type countingSink struct {
	counts map[s99logger.Level]int
}

func (s *countingSink) Write(_ context.Context, rec s99logger.Record) error {
	if s.counts == nil {
		s.counts = map[s99logger.Level]int{}
	}
	s.counts[rec.Level]++
	fmt.Printf("[%s] %s\n", rec.Level, rec.Message)
	return nil
}

func taskFinished(name string) s99logger.Event {
	return s99logger.NewEvent("task_finished", s99logger.String("name", name))
}

func main() {
	sink := &countingSink{}
	logger := s99logger.New(sink, s99logger.Options{})

	logger.Info(taskFinished("backup"))
	logger.Warn(taskFinished("cleanup"))
	logger.Error(taskFinished("report"))

	fmt.Printf("\ntotals: info=%d warn=%d error=%d\n",
		sink.counts[s99logger.LevelInfo],
		sink.counts[s99logger.LevelWarn],
		sink.counts[s99logger.LevelError],
	)
}
