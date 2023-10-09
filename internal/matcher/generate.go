package matcher

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	ErrMatchIndexTooLow  = errors.New("match index too low")
	ErrMatchIndexTooHigh = errors.New("match index too high")
)

func (matcher *Matcher) generate(cmd *cobra.Command, args []string) error {
	if matcher.matchIndex < 0 {
		return ErrMatchCountTooLow
	}
	if matcher.matchIndex > matcher.matchCount {
		return ErrMatchCountTooHigh
	}

	f, err := os.Create(matcher.messagesFilepath)
	if err != nil {
		return err
	}
	defer f.Close()

	for gifterUsername, options := range matcher.matches {
		gifter := matcher.participantsMap[gifterUsername]
		giftee := matcher.participantsMap[options[matcher.matchIndex]]
		details := giftee.DumpGifterMessage(matcher.eventName, gifter)
		details += "===========================================================\n\n"
		f.WriteString(details)
	}

	log.Info().Msg("generate executed")
	return nil
}
