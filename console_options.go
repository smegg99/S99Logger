package s99logger

import "strings"

const (
	// ANSIReset clears ANSI formatting.
	ANSIReset = "\x1b[0m"
	// ANSIFaint dims text in terminals that support ANSI styling.
	ANSIFaint = "\x1b[2m"

	ANSIBlack   = "\x1b[30m"
	ANSIRed     = "\x1b[31m"
	ANSIGreen   = "\x1b[32m"
	ANSIYellow  = "\x1b[33m"
	ANSIBlue    = "\x1b[34m"
	ANSIMagenta = "\x1b[35m"
	ANSICyan    = "\x1b[36m"
	ANSIWhite   = "\x1b[37m"

	ANSIBrightBlack   = "\x1b[90m"
	ANSIBrightRed     = "\x1b[91m"
	ANSIBrightGreen   = "\x1b[92m"
	ANSIBrightYellow  = "\x1b[93m"
	ANSIBrightBlue    = "\x1b[94m"
	ANSIBrightMagenta = "\x1b[95m"
	ANSIBrightCyan    = "\x1b[96m"
	ANSIBrightWhite   = "\x1b[97m"
)

var defaultLevelPrefixes = map[Level]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
}

var defaultLevelColors = map[Level]string{
	LevelDebug: ANSIBrightGreen,
	LevelInfo:  ANSIBrightCyan,
	LevelWarn:  ANSIBrightYellow,
	LevelError: ANSIBrightRed,
}

// WithColor forces color on or off, overriding terminal detection.
func (s *ConsoleSink) WithColor(on bool) *ConsoleSink {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.color = on
	return s
}

// WithLevelPrefix sets the default console prefix for level.
func (s *ConsoleSink) WithLevelPrefix(level Level, prefix string) *ConsoleSink {
	return s.WithLocalizedLevelPrefix("", level, prefix)
}

// WithLevelPrefixes sets default console prefixes. Missing levels keep their
// existing prefix.
func (s *ConsoleSink) WithLevelPrefixes(prefixes map[Level]string) *ConsoleSink {
	return s.WithLocalizedLevelPrefixes("", prefixes)
}

// WithLocalizedLevelPrefix sets the console prefix for level when a record uses
// lang. An empty lang configures the default prefix.
func (s *ConsoleSink) WithLocalizedLevelPrefix(lang string, level Level, prefix string) *ConsoleSink {
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

// WithLocalizedLevelPrefixes sets console prefixes for records using lang.
// Missing levels keep their existing prefix.
func (s *ConsoleSink) WithLocalizedLevelPrefixes(lang string, prefixes map[Level]string) *ConsoleSink {
	for level, prefix := range prefixes {
		s.WithLocalizedLevelPrefix(lang, level, prefix)
	}
	return s
}

// WithLevelColor sets the default ANSI color for level. Empty color disables
// coloring for that level.
func (s *ConsoleSink) WithLevelColor(level Level, color string) *ConsoleSink {
	return s.WithLocalizedLevelColor("", level, color)
}

// WithLevelColors sets default ANSI colors for levels. Missing levels keep their
// existing color.
func (s *ConsoleSink) WithLevelColors(colors map[Level]string) *ConsoleSink {
	return s.WithLocalizedLevelColors("", colors)
}

// WithLocalizedLevelColor sets the ANSI color for level when a record uses
// lang. An empty lang configures the default color.
func (s *ConsoleSink) WithLocalizedLevelColor(lang string, level Level, color string) *ConsoleSink {
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

// WithLocalizedLevelColors sets ANSI colors for records using lang. Missing
// levels keep their existing color.
func (s *ConsoleSink) WithLocalizedLevelColors(lang string, colors map[Level]string) *ConsoleSink {
	for level, color := range colors {
		s.WithLocalizedLevelColor(lang, level, color)
	}
	return s
}

// WithServiceColor sets the ANSI color used for the service prefix. Empty color
// disables service-prefix coloring.
func (s *ConsoleSink) WithServiceColor(color string) *ConsoleSink {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.serviceColor = color
	return s
}

// WithServicePrefixStyle sets the strings around the service prefix. The
// default style renders service as [service].
func (s *ConsoleSink) WithServicePrefixStyle(open, close string) *ConsoleSink {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.serviceOpen = open
	s.serviceClose = close
	return s
}

// WithFieldColor sets the ANSI color used for structured field keys. Empty
// color disables field-key coloring.
func (s *ConsoleSink) WithFieldColor(color string) *ConsoleSink {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fieldColor = color
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
