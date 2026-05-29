// examples/pluralization/main.go
package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/smegg99/s99logger"
	"github.com/smegg99/s99logger/i18n"
)

//go:embed locales/*.toml
var localesFS embed.FS

func jobsProcessed(count int) s99logger.Event {
	return s99logger.NewPluralEvent("jobs_processed", count, s99logger.Int("count", count))
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

func run() error {
	translator, err := i18n.New(i18n.Options{
		DefaultLanguage: "en",
		FS:              localesFS,
		Files:           []string{"locales/en.toml", "locales/pl.toml"},
		Decoders:        map[string]i18n.Decoder{"toml": toml.Unmarshal},
	})
	if err != nil {
		return err
	}

	logger := s99logger.New(s99logger.NewConsoleSink(os.Stdout), s99logger.Options{
		Language:   "en",
		Service:    "jobs",
		Translator: translator,
	})

	logger.Info(jobsProcessed(1))
	logger.Info(jobsProcessed(5), s99logger.LanguageContext("pl"))
	return nil
}
