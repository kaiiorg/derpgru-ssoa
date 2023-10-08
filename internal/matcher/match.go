package matcher

import (
	"math/rand"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func (matcher *Matcher) match(cmd *cobra.Command, args []string) error {
	log.Info().Msg("match executed")
	return nil
}

func (matcher *Matcher) matchShuffle() int {
	iterations := rand.Int31n(int32(len(matcher.participants) - 1)) + 1 // At least once, but no more than the number of participants minus 1
	log.Trace().
		Int32("iterations", iterations).
		Msg("shuffled participant list")
	for i := int32(0); i < iterations; i++ {
		rand.Shuffle(
			len(matcher.participants),
			func(i, j int) {
				matcher.participants[i], matcher.participants[j] = matcher.participants[j], matcher.participants[i]
			},
		)
	}
	return int(iterations)
}
