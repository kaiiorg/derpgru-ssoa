package matcher

import (
	"github.com/rs/zerolog/log"
	"reflect"

	"github.com/kaiiorg/derpgru-ssoa/internal/models/participant"

	"github.com/spf13/cobra"
)

func (matcher *Matcher) fields(cmd *cobra.Command, args []string) error {
	fields := reflect.VisibleFields(reflect.TypeOf(participant.Participant{}))

	for _, field := range fields {
		csvTag, found := field.Tag.Lookup("csv")
		if !found {
			continue
		}

		log.Info().Str("csv", csvTag).Str("field", field.Name).Send()
	}

	return nil
}
