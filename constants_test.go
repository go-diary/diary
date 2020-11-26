package diary

import "testing"

func TestLevelToText(t *testing.T) {
	tests := []struct{
		Level int
		Text string
	}{
		{
			Level: LevelUnknown,
			Text: TextLevelUnknown,
		},
		{
			Level: LevelTrace,
			Text: TextLevelTrace,
		},
		{
			Level: LevelDebug,
			Text: TextLevelDebug,
		},
		{
			Level: LevelInfo,
			Text: TextLevelInfo,
		},
		{
			Level: LevelNotice,
			Text: TextLevelNotice,
		},
		{
			Level: LevelWarning,
			Text: TextLevelWarning,
		},
		{
			Level: LevelError,
			Text: TextLevelError,
		},
		{
			Level: LevelFatal,
			Text: TextLevelFatal,
		},
		{
			Level: 99,
			Text: TextLevelUnknown,
		},
	}

	for i, test := range tests {
		text := levelToText(test.Level)
		if test.Text != text {
			t.Errorf("[%d] | level: %d | text = %v; want %v", i, test.Level, text, test.Text)
		}
	}
}

func TestTextToLevel(t *testing.T) {
	tests := []struct{
		Level int
		Text string
	}{
		{
			Level: LevelUnknown,
			Text: TextLevelUnknown,
		},
		{
			Level: LevelTrace,
			Text: TextLevelTrace,
		},
		{
			Level: LevelDebug,
			Text: TextLevelDebug,
		},
		{
			Level: LevelInfo,
			Text: TextLevelInfo,
		},
		{
			Level: LevelNotice,
			Text: TextLevelNotice,
		},
		{
			Level: LevelWarning,
			Text: TextLevelWarning,
		},
		{
			Level: LevelError,
			Text: TextLevelError,
		},
		{
			Level: LevelFatal,
			Text: TextLevelFatal,
		},
		{
			Level: LevelUnknown,
			Text: "something",
		},
	}

	for i, test := range tests {
		level := textToLevel(test.Text)
		if test.Level != level {
			t.Errorf("[%d] | text: %s | level = %v; want %v", i, test.Text, level, test.Level)
		}
	}
}