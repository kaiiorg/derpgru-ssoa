package matcher

import (
	"github.com/kaiiorg/derpgru-ssoa/internal/models/participant"

	"github.com/spf13/cobra"
)

const (
	MATCH_CMD_NAME    = "match"
	GENERATE_CMD_NAME = "generate"

	LOG_LEVEL_FLAG = "log-level"
	PARTICIPANTS_INPUT_FLAG = "in"
	MATCHES_FLAG = "matches"
	MESSAGES_FLAG = "messages"
	EVENT_NAME_FLAG = "name"
	INDEX_SELECT_FLAG = "select"
)

type Matcher struct {
	logLevel string

	participantsFilepath string
	matchesFilepath string
	messagesFilepath string

	participants []*participant.Participant

	// Used with generate command
	eventName string
	matchIndex int

	rootCmd *cobra.Command
	matchCmd *cobra.Command
	generateCmd *cobra.Command
}

func New() *Matcher {
	return &Matcher{}
}

func (matcher *Matcher) CobraCommand() *cobra.Command {
	matcher.rootCmd = &cobra.Command{
		Use:               "matcher",
		Short:             "Derpgru Secret Something or Another Matcher",
		Long:              `Program for matching Secret Santa participants and generating messages once matches have been vetted by organizers.`,
		PersistentPreRunE: matcher.configure,
	}
	matcher.rootCmd.PersistentFlags().StringVarP(&matcher.participantsFilepath, PARTICIPANTS_INPUT_FLAG, "i", "./pii/in.csv", "Input CSV file")
	matcher.rootCmd.PersistentFlags().StringVarP(&matcher.matchesFilepath, MATCHES_FLAG, "m", "./pii/matches.csv", "Matches output/input CSV file")
	matcher.rootCmd.PersistentFlags().StringVarP(&matcher.messagesFilepath, MESSAGES_FLAG, "o", "./pii/messages.txt", "Messages generated from selected matches file")
	matcher.rootCmd.PersistentFlags().StringVarP(&matcher.logLevel, LOG_LEVEL_FLAG, "l", "info", "zerolog.LogLevel; defaults to info if an invalid value is provided")

	matcher.matchCmd = &cobra.Command{
		Use:     MATCH_CMD_NAME,
		Short:   "Generate multiple matches for all participants and dump them to file",
		PreRunE: matcher.load,
		RunE:    matcher.match,
	}
	matcher.rootCmd.AddCommand(matcher.matchCmd)

	matcher.generateCmd = &cobra.Command{
		Use:     GENERATE_CMD_NAME,
		Short:   "Generate messages from given match list",
		PreRunE: matcher.load,
		RunE:    matcher.generate,
	}
	matcher.generateCmd.Flags().StringVarP(&matcher.eventName, EVENT_NAME_FLAG, "n", "", "Event name")
	matcher.generateCmd.Flags().IntVarP(&matcher.matchIndex, INDEX_SELECT_FLAG, "s", 0, "Which match pair to use")
	matcher.generateCmd.MarkFlagRequired(EVENT_NAME_FLAG)
	matcher.rootCmd.AddCommand(matcher.generateCmd)

	return matcher.rootCmd
}
