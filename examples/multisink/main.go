// examples/multisink/main.go
package main

import (
	"os"

	"github.com/smegg99/s99logger"
)

func requestHandled(path string) s99logger.Event {
	return s99logger.NewEvent("request_handled", s99logger.String("path", path))
}

func main() {
	sink := s99logger.MultiSink(
		s99logger.NewConsoleSink(os.Stdout),
		s99logger.NewJSONSink(os.Stdout),
	)
	logger := s99logger.New(sink, s99logger.Options{Service: "api"})

	logger.Info(requestHandled("/health"))
}
