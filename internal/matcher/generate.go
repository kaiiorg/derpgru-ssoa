package matcher

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func (matcher *Matcher) generate(cmd *cobra.Command, args []string) error {
	log.Info().Msg("generate executed")
	return nil
}
