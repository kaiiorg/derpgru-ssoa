package matcher

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func (matcher *Matcher) match(cmd *cobra.Command, args []string) error {
	log.Info().Msg("match executed")
	return nil
}
