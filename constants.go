package diary

import "strings"

const (
	LevelTrace = 0
	LevelDebug = 1
	LevelInfo = 2
	LevelNotice = 3
	LevelWarning = 4
	LevelError = 5
	LevelFatal = 6
)

const (
	TextLevelTrace = "trace"
	TextLevelTraceEnter = "enter"
	TextLevelTraceExit = "exit"
	TextLevelDebug = "debug"
	TextLevelInfo = "info"
	TextLevelNotice = "notice"
	TextLevelWarning = "warning"
	TextLevelError = "error"
	TextLevelFatal = "fatal"
)

func ConvertFromTextLevel(value string) int {
	switch strings.ToLower(value) {
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
	return -1
}

func IsValidLevel(value int) bool {
	switch value {
	case LevelTrace:
	case LevelDebug:
	case LevelInfo:
	case LevelNotice:
	case LevelWarning:
	case LevelError:
	case LevelFatal:
		return true
	}
	return false
}