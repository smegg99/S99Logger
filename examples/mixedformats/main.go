// examples/mixedformats/main.go

package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/smegg99/s99logger"
	"github.com/smegg99/s99logger/examples/internal/sampleevents"
	"github.com/smegg99/s99logger/i18n"
	"go.yaml.in/yaml/v3"
)

//go:embed locales/*.json locales/*.toml locales/*.yaml
var localesFS embed.FS

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
		Files: []string{
			"locales/en.json",
			"locales/pl.toml",
			"locales/es.yaml",
		},
		Decoders: map[string]i18n.Decoder{
			"toml": toml.Unmarshal,
			"yaml": yaml.Unmarshal,
		},
	})
	if err != nil {
		return fmt.Errorf("setup translator: %w", err)
	}

	logger := s99logger.New(
		s99logger.NewConsoleSink(os.Stdout),
		s99logger.Options{Language: "en", Service: "formats", Translator: translator},
	)

	logger.Info(sampleevents.AppStarted("formats", "1.0.0"))
	logger.Warn(sampleevents.ReconnectAttempt("cache-1", 1, 3, 250*time.Millisecond), s99logger.LanguageContext("pl"))
	logger.Error(sampleevents.LoginFailed("alice", errors.New("invalid password")), s99logger.LanguageContext("es"))
	logger.Info(sampleevents.JobsProcessed(2), s99logger.LanguageContext("es"))

	return nil
}
