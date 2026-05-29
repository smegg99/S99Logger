// examples/minimal/main.go
//
// The smallest possible setup: a colored console logger with no translator.
// Without a translator, the message is simply the stable message id.

package main

import (
	"os"

	"github.com/smegg99/s99logger"
)

func main() {
	logger := s99logger.New(
		s99logger.NewConsoleSink(os.Stderr),
		s99logger.Options{Service: "minimal"},
	)

	logger.Info(s99logger.NewEvent("server_starting", s99logger.Int("port", 8080)))
}
