// examples/scopedattrs/main.go
package main

import (
	"os"

	"github.com/smegg99/s99logger"
)

func cacheMiss(key string) s99logger.Event {
	return s99logger.NewEvent("cache_miss", s99logger.String("key", key))
}

func main() {
	logger := s99logger.New(s99logger.NewConsoleSink(os.Stdout), s99logger.Options{
		Service:  "cache",
		MinLevel: s99logger.LevelInfo,
	})
	requestLogger := logger.With(s99logger.String("request_id", "req-42"))

	requestLogger.Debug(cacheMiss("users:1")) // filtered by MinLevel
	requestLogger.Info(cacheMiss("users:1"))
}
