// examples/localizedprefixes/main.go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/smegg99/s99logger"
)

type staticTranslator map[string]map[s99logger.MessageID]string

func (t staticTranslator) Translate(_ context.Context, lang string, event s99logger.Event) (string, error) {
	if messages := t[lang]; messages != nil {
		if msg, ok := messages[event.MessageID()]; ok {
			return msg, nil
		}
	}
	if messages := t["en"]; messages != nil {
		if msg, ok := messages[event.MessageID()]; ok {
			return msg, nil
		}
	}
	return "", fmt.Errorf("missing message %q", event.MessageID())
}

func serviceReady(name string) s99logger.Event {
	return s99logger.NewEvent("service_ready", s99logger.String("name", name))
}

func serviceError(name string) s99logger.Event {
	return s99logger.NewEvent("service_error", s99logger.String("name", name))
}

func main() {
	sink := s99logger.NewConsoleSink(os.Stdout).
		WithLocalizedLevelPrefixes("en", map[s99logger.Level]string{
			s99logger.LevelInfo:  "INFO",
			s99logger.LevelError: "ERROR",
		}).
		WithLocalizedLevelPrefixes("pl", map[s99logger.Level]string{
			s99logger.LevelInfo:  "OKEJ",
			s99logger.LevelError: "BLAD",
		}).
		WithLocalizedLevelColors("pl", map[s99logger.Level]string{
			s99logger.LevelInfo:  s99logger.ANSIGreen,
			s99logger.LevelError: s99logger.ANSIRed,
		})

	logger := s99logger.New(sink, s99logger.Options{
		Language: "en",
		Service:  "localized-prefixes",
		Translator: staticTranslator{
			"en": {
				"service_ready": "Service is ready",
				"service_error": "Service failed",
			},
			"pl": {
				"service_ready": "Usluga jest gotowa",
				"service_error": "Usluga nie dziala",
			},
		},
	})

	logger.Info(serviceReady("api"))

	ctx := s99logger.LanguageContext("pl")
	logger.Info(serviceReady("api"), ctx)
	logger.Error(serviceError("api"), ctx)
}
