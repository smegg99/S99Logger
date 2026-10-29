// i18n/verify.go

package i18n

import (
	"fmt"
	"strings"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/smegg99/s99logger"
)

// Verify lists every (lang, id) with no translation, without falling back; call it from a test so a missing one fails the build.
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
