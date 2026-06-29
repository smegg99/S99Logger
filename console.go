// console.go

package s99logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// ConsoleSink writes human-readable, optionally colored lines, suited for the
// terminal. Color is auto-enabled when the writer is a terminal. Level prefixes
// and colors can be customized, including per-language prefixes selected from
// each record's language. It is safe for concurrent use.
type ConsoleSink struct {
	mu            sync.Mutex
	out           io.Writer
	color         bool
	timeFormat    string
	levelPrefixes map[string]map[Level]string
	levelColors   map[string]map[Level]string
	serviceColor  string
	serviceOpen   string
	serviceClose  string
	fieldColor    string
}

// NewConsoleSink returns a ConsoleSink writing to w. Color is enabled only when
// w is a terminal.
func NewConsoleSink(w io.Writer) *ConsoleSink {
	if w == nil {
		w = io.Discard
	}
	return &ConsoleSink{
		out:           w,
		color:         isTerminal(w),
		timeFormat:    time.DateTime,
		levelPrefixes: map[string]map[Level]string{"": cloneLevelStrings(defaultLevelPrefixes)},
		levelColors:   map[string]map[Level]string{"": cloneLevelStrings(defaultLevelColors)},
		serviceColor:  ANSIFaint,
		serviceOpen:   "[",
		serviceClose:  "]",
		fieldColor:    ANSIFaint,
	}
}

// Write renders rec as a single human-readable line.
func (s *ConsoleSink) Write(_ context.Context, rec Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var b strings.Builder

	b.WriteString(rec.Time.Local().Format(s.timeFormat))
	b.WriteByte(' ')
	if rec.Service != "" {
		s.writeService(&b, rec.Service)
		b.WriteByte(' ')
	}
	s.writeLevel(&b, rec.Lang, rec.Level)
	b.WriteByte(' ')
	b.WriteString(rec.Message)

	// Always surface the stable message id. Skip it when it equals the message
	// (i.e. there was no translation) to avoid redundant output.
	if rec.Message != string(rec.MessageID) {
		s.writeField(&b, "id", string(rec.MessageID))
	}
	for _, a := range rec.Attrs {
		// Skip attrs that collide with fields already printed above, matching
		// the JSON sink's reserved-key protection.
		if _, ok := reservedKeys[a.Key]; ok || a.Key == "id" {
			continue
		}
		s.writeField(&b, a.Key, fmt.Sprintf("%v", a.Value))
	}
	b.WriteByte('\n')

	if _, err := io.WriteString(s.out, b.String()); err != nil {
		return fmt.Errorf("%w: %w", errWriteConsoleRecord, err)
	}
	return nil
}

// writeLevel writes the configured level prefix, padded to a fixed width and
// colored when color is enabled.
func (s *ConsoleSink) writeLevel(b *strings.Builder, lang string, level Level) {
	prefix := padLevelPrefix(s.levelPrefix(lang, level))
	color := s.levelColor(lang, level)
	if s.color && color != "" {
		b.WriteString(color)
		b.WriteString(prefix)
		b.WriteString(ANSIReset)
		return
	}
	b.WriteString(prefix)
}

// writeService writes the service as a prefix segment instead of a structured
// field, keeping the console line scannable while JSON keeps service as data.
func (s *ConsoleSink) writeService(b *strings.Builder, service string) {
	if s.color && s.serviceColor != "" {
		b.WriteString(s.serviceColor)
		b.WriteString(s.serviceOpen)
		b.WriteString(service)
		b.WriteString(s.serviceClose)
		b.WriteString(ANSIReset)
		return
	}
	b.WriteString(s.serviceOpen)
	b.WriteString(service)
	b.WriteString(s.serviceClose)
}

func (s *ConsoleSink) levelPrefix(lang string, level Level) string {
	s.ensureLevelPrefixes()
	for _, key := range consoleLanguageKeys(lang) {
		if prefixes := s.levelPrefixes[key]; prefixes != nil {
			if prefix, ok := prefixes[level]; ok {
				return prefix
			}
		}
	}
	return strings.ToUpper(level.String())
}

func (s *ConsoleSink) levelColor(lang string, level Level) string {
	s.ensureLevelColors()
	for _, key := range consoleLanguageKeys(lang) {
		if colors := s.levelColors[key]; colors != nil {
			if color, ok := colors[level]; ok {
				return color
			}
		}
	}
	return ""
}

func consoleLanguageKeys(lang string) []string {
	key := normalizeConsoleLanguage(lang)
	if key == "" {
		return []string{""}
	}
	keys := []string{key}
	if i := strings.Index(key, "-"); i > 0 {
		keys = append(keys, key[:i])
	}
	return append(keys, "")
}

func padLevelPrefix(prefix string) string {
	if width := utf8.RuneCountInString(prefix); width < 5 {
		return prefix + strings.Repeat(" ", 5-width)
	}
	return prefix
}

// writeField appends a " key=value" pair, quoting values that contain spaces
// and dimming the key when color is enabled.
func (s *ConsoleSink) writeField(b *strings.Builder, key, value string) {
	if strings.ContainsAny(value, " \t\r\n\"") {
		value = strconv.Quote(value)
	}
	b.WriteByte(' ')
	if s.color && s.fieldColor != "" {
		b.WriteString(s.fieldColor)
		b.WriteString(key)
		b.WriteByte('=')
		b.WriteString(ANSIReset)
		b.WriteString(value)
		return
	}
	b.WriteString(key)
	b.WriteByte('=')
	b.WriteString(value)
}

// isTerminal reports whether w is a character device (a terminal).
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
