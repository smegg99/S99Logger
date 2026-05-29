// examples/consoleprefixes/main.go
package main

import (
	"os"

	"github.com/smegg99/s99logger"
)

func taskDone(name string) s99logger.Event {
	return s99logger.NewEvent("task_done", s99logger.String("task", name))
}

func main() {
	sink := s99logger.NewConsoleSink(os.Stdout).
		WithLevelPrefixes(map[s99logger.Level]string{
			s99logger.LevelInfo:  "OK",
			s99logger.LevelWarn:  "WAIT",
			s99logger.LevelError: "FAIL",
		}).
		WithLevelColors(map[s99logger.Level]string{
			s99logger.LevelInfo:  s99logger.ANSIGreen,
			s99logger.LevelWarn:  s99logger.ANSIYellow,
			s99logger.LevelError: s99logger.ANSIRed,
		}).
		WithServicePrefixStyle("{", "}").
		WithServiceColor(s99logger.ANSIBrightWhite)

	logger := s99logger.New(sink, s99logger.Options{Service: "console"})

	logger.Info(taskDone("backup"))
	logger.Warn(taskDone("cleanup"))
}
