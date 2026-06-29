// context.go

package s99logger

import "context"

type langKey struct{}

// ContextWithLanguage returns a copy of ctx carrying lang. When a Logger logs
// with this context, lang overrides Options.Language for that call, which makes
// per-request or per-user localization possible without a logger per language.
func ContextWithLanguage(ctx context.Context, lang string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, langKey{}, lang)
}

// LanguageContext returns a background context carrying lang.
func LanguageContext(lang string) context.Context {
	return ContextWithLanguage(context.Background(), lang)
}

// LanguageFromContext returns the language stored by ContextWithLanguage, and
// whether one was present (and non-empty).
func LanguageFromContext(ctx context.Context) (string, bool) {
	lang, ok := ctx.Value(langKey{}).(string)
	return lang, ok && lang != ""
}
