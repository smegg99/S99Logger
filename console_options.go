// console_options.go

package s99logger

import "strings"

const (
	// ansiReset clears ANSI formatting.
	ansiReset = "\x1b[0m"
	// ansiFaint dims text in terminals that support ANSI styling.
	ansiFaint = "\x1b[2m"

	ANSIRed          = "\x1b[31m"
	ANSIGreen        = "\x1b[32m"
	ANSIYellow       = "\x1b[33m"
	ansiBrightRed    = "\x1b[91m"
	ansiBrightGreen  = "\x1b[92m"
	ansiBrightYellow = "\x1b[93m"
	ansiBrightCyan   = "\x1b[96m"
	ANSIBrightWhite  = "\x1b[97m"
)

var defaultLevelPrefixes = map[Level]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
}

var defaultLevelColors = map[Level]string{
	LevelDebug: ansiBrightGreen,
	LevelInfo:  ansiBrightCyan,
	LevelWarn:  ansiBrightYellow,
	LevelError: ansiBrightRed,
}

// WithColor forces color on or off, overriding terminal detection.
func (s *ConsoleSink) WithColor(on bool) *ConsoleSink {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.color = on
	return s
}

// WithLevelPrefixes sets default console prefixes; missing levels keep theirs.
func (s *ConsoleSink) WithLevelPrefixes(prefixes map[Level]string) *ConsoleSink {
	return s.WithLocalizedLevelPrefixes("", prefixes)
}

// withLocalizedLevelPrefix sets the prefix for level when a record uses lang; an empty lang sets the default.
func (s *ConsoleSink) withLocalizedLevelPrefix(lang string, level Level, prefix string) *ConsoleSink {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureLevelPrefixes()
	key := normalizeConsoleLanguage(lang)
	if s.levelPrefixes[key] == nil {
		s.levelPrefixes[key] = map[Level]string{}
	}
	s.levelPrefixes[key][level] = prefix
	return s
}

// WithLocalizedLevelPrefixes sets console prefixes for records in lang; missing levels keep theirs.
func (s *ConsoleSink) WithLocalizedLevelPrefixes(lang string, prefixes map[Level]string) *ConsoleSink {
	for level, prefix := range prefixes {
		s.withLocalizedLevelPrefix(lang, level, prefix)
	}
	return s
}

// WithLevelColors sets default ANSI colors for levels; missing levels keep theirs.
func (s *ConsoleSink) WithLevelColors(colors map[Level]string) *ConsoleSink {
	return s.WithLocalizedLevelColors("", colors)
}

// withLocalizedLevelColor sets the color for level when a record uses lang; an empty lang sets the default.
func (s *ConsoleSink) withLocalizedLevelColor(lang string, level Level, color string) *ConsoleSink {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureLevelColors()
	key := normalizeConsoleLanguage(lang)
	if s.levelColors[key] == nil {
		s.levelColors[key] = map[Level]string{}
	}
	s.levelColors[key][level] = color
	return s
}

// WithLocalizedLevelColors sets ANSI colors for records in lang; missing levels keep theirs.
func (s *ConsoleSink) WithLocalizedLevelColors(lang string, colors map[Level]string) *ConsoleSink {
	for level, color := range colors {
		s.withLocalizedLevelColor(lang, level, color)
	}
	return s
}

// WithServiceColor sets the service prefix's ANSI color; empty disables it.
func (s *ConsoleSink) WithServiceColor(color string) *ConsoleSink {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.serviceColor = color
	return s
}

// WithServicePrefixStyle sets the strings around the service prefix, [service] by default.
func (s *ConsoleSink) WithServicePrefixStyle(open, close string) *ConsoleSink {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.serviceOpen = open
	s.serviceClose = close
	return s
}

func (s *ConsoleSink) ensureLevelPrefixes() {
	if s.levelPrefixes == nil {
		s.levelPrefixes = map[string]map[Level]string{
			"": cloneLevelStrings(defaultLevelPrefixes),
		}
	}
}

func (s *ConsoleSink) ensureLevelColors() {
	if s.levelColors == nil {
		s.levelColors = map[string]map[Level]string{
			"": cloneLevelStrings(defaultLevelColors),
		}
	}
}

func cloneLevelStrings(input map[Level]string) map[Level]string {
	output := make(map[Level]string, len(input))
	for level, value := range input {
		output[level] = value
	}
	return output
}

func normalizeConsoleLanguage(lang string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(lang), "_", "-"))
}
