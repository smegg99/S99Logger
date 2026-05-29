// examples/jsonstdout/main.go
//
// Newline-delimited JSON to stdout, no translator. Pipe it into jq or a log
// shipper. Each line is one self-contained JSON object.

package main

import (
	"os"

	"github.com/smegg99/s99logger"
)

func requestHandled(method, path string, status int) s99logger.Event {
	return s99logger.NewEvent("request_handled",
		s99logger.String("method", method),
		s99logger.String("path", path),
		s99logger.Int("status", status),
	)
}

func main() {
	logger := s99logger.New(
		s99logger.NewJSONSink(os.Stdout),
		s99logger.Options{Service: "api"},
	)
	logger.Info(requestHandled("GET", "/health", 200))
	logger.Warn(requestHandled("POST", "/login", 429))
	logger.Error(requestHandled("GET", "/missing", 404))
}
