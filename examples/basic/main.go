// examples/basic/main.go

package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/smegg99/s99logger"
	"github.com/smegg99/s99logger/examples/internal/sampleevents"
	"github.com/smegg99/s99logger/i18n"
)

// The application owns its locale files and embeds them. The fs.FS contract
// means this works equally well with os.DirFS during development.
//
//go:embed locales/*.toml
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
		Files:           []string{"locales/en.toml", "locales/pl.toml"},
		Decoders:        map[string]i18n.Decoder{"toml": toml.Unmarshal},
	})
	if err != nil {
		return fmt.Errorf("setup translator: %w", err)
	}

	logger := s99logger.New(s99logger.NewConsoleSink(os.Stderr), s99logger.Options{
		Language:   "pl",
		Service:    "demo",
		Translator: translator,
	})

	logger.Info(sampleevents.AppStarted("demo", "1.4.2"))
	return nil
}
