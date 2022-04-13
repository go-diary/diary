package diary

import "strings"

const (
	LevelTrace   = 0
	LevelDebug   = 1
	LevelInfo    = 2
	LevelNotice  = 3
	LevelWarning = 4
	LevelError   = 5
	LevelFatal   = 6
	LevelAudit   = 7
)

var Levels = []int{
	LevelTrace,
	LevelDebug,
	LevelInfo,
	LevelNotice,
	LevelWarning,
	LevelWarning,
	LevelError,
	LevelFatal,
	LevelAudit,
}

const (
	TextLevelTrace      = "trace"
	TextLevelTraceEnter = "enter"
	TextLevelTraceExit  = "exit"
	TextLevelDebug      = "debug"
	TextLevelInfo       = "info"
	TextLevelNotice     = "notice"
	TextLevelWarning    = "warning"
	TextLevelError      = "error"
	TextLevelFatal      = "fatal"
	TextLevelAudit      = "audit"
)

var TextLevels = []string{
	TextLevelTrace,
	TextLevelDebug,
	TextLevelInfo,
	TextLevelNotice,
	TextLevelWarning,
	TextLevelWarning,
	TextLevelError,
	TextLevelFatal,
	TextLevelAudit,
}

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
	case TextLevelAudit:
		return LevelAudit
	}
	return -1
}

func IsValidLevel(value int) bool {
	switch value {
	case LevelTrace:
		return true
	case LevelDebug:
		return true
	case LevelInfo:
		return true
	case LevelNotice:
		return true
	case LevelWarning:
		return true
	case LevelError:
		return true
	case LevelFatal:
		return true
	case LevelAudit:
		return true
	}
	return false
}
