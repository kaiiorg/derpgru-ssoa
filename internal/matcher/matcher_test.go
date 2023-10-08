package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/spf13/cobra"
)

func testingCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "testing",
		Short:             "Testing command",
	}
}

func TestNew(t *testing.T) {
	// Act
	matcher := New()

	// Assert
	require.NotNil(t, matcher)
}

func TestCobraCommand(t *testing.T) {
	// Arrange
	matcher := New()

	// Act
	cmd := matcher.CobraCommand()

	// Assert
	require.NotNil(t, cmd)

	// rootCmd
	require.NotNil(t, matcher.rootCmd)
	require.NotNil(t, matcher.rootCmd.PersistentPreRunE)
	rootFlags := matcher.rootCmd.PersistentFlags()
	require.NotNilf(t, rootFlags.Lookup(PARTICIPANTS_INPUT_FLAG), "Did not find root persistent flag '%s'", PARTICIPANTS_INPUT_FLAG)
	require.NotNilf(t, rootFlags.Lookup(MATCHES_FLAG), "Did not find root persistent flag '%s'", MATCHES_FLAG)
	require.NotNilf(t, rootFlags.Lookup(MESSAGES_FLAG), "Did not find root persistent flag '%s'", MESSAGES_FLAG)
	require.NotNilf(t, rootFlags.Lookup(LOG_LEVEL_FLAG), "Did not find root persistent flag '%s'", LOG_LEVEL_FLAG)
	for _, command := range matcher.rootCmd.Commands() {
		// Make sure each of our expected commands got added to the root command
		require.Contains(t, []string{MATCH_CMD_NAME, GENERATE_CMD_NAME}, command.Use)
	}

	// matchCmd
	require.NotNil(t, matcher.matchCmd)
	require.NotNil(t, matcher.matchCmd.PreRunE)

	// generateCmd
	require.NotNil(t, matcher.generateCmd)
	require.NotNil(t, matcher.generateCmd.PreRunE)
	generateFlags := matcher.generateCmd.Flags()
	require.NotNilf(t, generateFlags.Lookup(EVENT_NAME_FLAG), "Did not find generate flag '%s'", EVENT_NAME_FLAG)
	require.NotNilf(t, generateFlags.Lookup(INDEX_SELECT_FLAG), "Did not find generate flag '%s'", INDEX_SELECT_FLAG)
}
