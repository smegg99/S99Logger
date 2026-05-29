package s99logger

import (
	"context"
	"time"
)

// Debug logs an event at debug level.
func (l *Logger) Debug(event Event, ctx ...context.Context) { l.log(LevelDebug, event, ctx...) }

// Info logs an event at info level.
func (l *Logger) Info(event Event, ctx ...context.Context) { l.log(LevelInfo, event, ctx...) }

// Warn logs an event at warn level.
func (l *Logger) Warn(event Event, ctx ...context.Context) { l.log(LevelWarn, event, ctx...) }

// Error logs an event at error level.
func (l *Logger) Error(event Event, ctx ...context.Context) { l.log(LevelError, event, ctx...) }

// log resolves the event into a Record and writes it to the sink. It never
// fails: translation errors fall back to the message id and sink errors are
// discarded.
func (l *Logger) log(level Level, event Event, ctxs ...context.Context) {
	if l == nil || isNilInterface(event) || isNilInterface(l.sink) {
		return
	}
	if level < l.minLevel {
		return
	}

	ctx := context.Background()
	if len(ctxs) > 0 && ctxs[0] != nil {
		ctx = ctxs[0]
	}

	lang := l.lang
	if cl, ok := LanguageFromContext(ctx); ok {
		lang = cl
	}

	clock := l.clock
	if clock == nil {
		clock = time.Now
	}

	rec := Record{
		Time:      clock(),
		Level:     level,
		Service:   l.service,
		Lang:      lang,
		MessageID: event.MessageID(),
		Message:   string(event.MessageID()),
		Attrs:     l.mergeAttrs(event.Attrs()),
	}
	if l.translator != nil {
		if msg, err := l.translator.Translate(ctx, lang, event); err == nil {
			rec.Message = msg
		}
	}
	_ = l.sink.Write(ctx, rec)
}
