// examples/rotatingfile/main.go
package main

import (
	"fmt"

	"github.com/smegg99/s99logger"
	"github.com/smegg99/s99logger/rotation"
)

func cleanupFinished(files int) s99logger.Event {
	return s99logger.NewEvent("cleanup_finished", s99logger.Int("files", files))
}

func main() {
	sink, err := rotation.New(rotation.Options{
		Directory:  "logs",
		Filename:   "app.log",
		LocalTime:  true,
		MaxBackups: 7,
	})
	if err != nil {
		panic(err)
	}
	defer sink.Close()

	logger := s99logger.New(sink, s99logger.Options{Service: "worker"})
	logger.Info(cleanupFinished(12))

	fmt.Println("wrote logs/app.log")
}
