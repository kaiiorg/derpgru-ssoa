package matcher

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestConfigure_NoError(t *testing.T) {
	// Arrange
	matcher := New()

	// Act
	err := matcher.configure(testingCmd(), []string{})

	// Assert
	require.NoError(t, err)
}

func TestConfigureLogLevel(t *testing.T) {
	// Arrange
	matcher := New()
	testLevels := map[string]zerolog.Level{
		"":         zerolog.InfoLevel,
		"invalid":  zerolog.InfoLevel,
		"none":     zerolog.InfoLevel,
		"trace":    zerolog.TraceLevel,
		"debug":    zerolog.DebugLevel,
		"info":     zerolog.InfoLevel,
		"warn":     zerolog.WarnLevel,
		"error":    zerolog.ErrorLevel,
		"disabled": zerolog.Disabled,
	}

	// Act
	for testValue, expectedLevel := range testLevels {
		matcher.logLevel = testValue
		matcher.configureLogLevel()

		// Assert
		require.Equal(t, expectedLevel, zerolog.GlobalLevel())
	}
}
