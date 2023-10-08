package matcher

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func (matcher *Matcher) load(cmd *cobra.Command, args []string) error {
	// Load participants for all commands
	participantsF, err := os.Open(matcher.participantsFilepath)
	if err != nil {
		return err
	}
	defer participantsF.Close()

	err = gocsv.UnmarshalFile(participantsF, &matcher.participants)
	if err != nil {
		return err
	}

	log.Info().Str("participantsFilepath", matcher.participantsFilepath).Msg("loaded participants file")

	return nil

	// TODO Load the matches
	/*
	// We're done if this isn't the generate command
	if cmd.Use != GENERATE_CMD_NAME {
		return nil
	}

	return matcher.loadMatches()
	*/
}

/*
func (matcher *Matcher) loadMatches() error {
	// TODO Load the matches
	return nil
}
*/