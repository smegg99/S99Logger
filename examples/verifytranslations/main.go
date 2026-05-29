// examples/verifytranslations/main.go
package main

import (
	"fmt"
	"testing/fstest"

	"github.com/smegg99/s99logger"
	"github.com/smegg99/s99logger/i18n"
)

func main() {
	locales := fstest.MapFS{
		"locales/en.json": {Data: []byte(`{"app_started":"App started"}`)},
		"locales/pl.json": {Data: []byte(`{"app_started":"Aplikacja uruchomiona"}`)},
	}
	translator, err := i18n.New(i18n.Options{
		DefaultLanguage: "en",
		FS:              locales,
		Files:           []string{"locales/en.json", "locales/pl.json"},
	})
	if err != nil {
		panic(err)
	}

	if err := translator.Verify([]string{"en", "pl"}, s99logger.MessageID("app_started")); err != nil {
		panic(err)
	}

	fmt.Println("translations complete")
}
