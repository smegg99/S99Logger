package i18n

import (
	"fmt"
	"strings"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/smegg99/s99logger"
)

// Verify reports every (lang, id) combination that has no translation, checking
// each language independently without falling back to the default. Call it from
// a test with your event ids so a missing or mistyped translation fails the
// build instead of silently degrading to the message id in production.
func (t *Translator) Verify(langs []string, ids ...s99logger.MessageID) error {
	var missing []string
	for _, lang := range langs {
		localizer := goi18n.NewLocalizer(t.bundle, lang)
		for _, id := range ids {
			if _, err := localizer.Localize(&goi18n.LocalizeConfig{MessageID: string(id)}); err != nil {
				missing = append(missing, fmt.Sprintf("%s/%s", lang, id))
			}
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("%w: %s", errMissingTranslations, strings.Join(missing, ", "))
	}
	return nil
}
