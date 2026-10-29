// context.go

package s99logger

import "context"

type langKey struct{}

// contextWithLanguage returns ctx carrying lang, which overrides Options.Language for calls logged with it.
func contextWithLanguage(ctx context.Context, lang string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, langKey{}, lang)
}

// LanguageContext returns a background context carrying lang.
func LanguageContext(lang string) context.Context {
	return contextWithLanguage(context.Background(), lang)
}

// languageFromContext returns the non-empty language stored by contextWithLanguage.
func languageFromContext(ctx context.Context) (string, bool) {
	lang, ok := ctx.Value(langKey{}).(string)
	return lang, ok && lang != ""
}
