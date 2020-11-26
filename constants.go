package diary

import "strings"

const (
	LevelUnknown = -1
	LevelTrace = 0
	LevelDebug = 1
	LevelInfo = 2
	LevelNotice = 3
	LevelWarning = 4
	LevelError = 5
	LevelFatal = 6
)

const (
	TextLevelUnknown = "unknown"
	TextLevelTrace = "trace"
	TextLevelDebug = "debug"
	TextLevelInfo = "info"
	TextLevelNotice = "notice"
	TextLevelWarning = "warning"
	TextLevelError = "error"
	TextLevelFatal = "fatal"
)

var levelToText = func(level int) string {
	switch level {
	case LevelTrace:
		return TextLevelTrace
	case LevelDebug:
		return TextLevelDebug
	case LevelInfo:
		return TextLevelInfo
	case LevelNotice:
		return TextLevelNotice
	case LevelWarning:
		return TextLevelWarning
	case LevelError:
		return TextLevelError
	case LevelFatal:
		return TextLevelFatal
	}
	return TextLevelUnknown
}

var textToLevel = func(text string) int {
	switch strings.ToLower(strings.TrimSpace(text)) {
	case TextLevelTrace:
		return LevelTrace
	case TextLevelDebug:
		return LevelDebug
	case TextLevelInfo:
		return LevelInfo
	case TextLevelNotice:
		return LevelNotice
	case TextLevelWarning:
		return LevelWarning
	case TextLevelError:
		return LevelError
	case TextLevelFatal:
		return LevelFatal
	}
	return LevelUnknown
}