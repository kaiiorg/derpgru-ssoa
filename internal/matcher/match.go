package matcher

import (
	"errors"
	"math/rand"

	"github.com/kaiiorg/derpgru-ssoa/internal/models/participant"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	ErrMatchCountTooLow  = errors.New("match count too low")
	ErrMatchCountTooHigh = errors.New("match count too high")
)

func (matcher *Matcher) match(cmd *cobra.Command, args []string) error {
	if matcher.matchCount < 1 {
		log.Error().Int("requestedMatchCount", matcher.matchCount).Msg("Requested match count is too low!")
		return ErrMatchCountTooLow
	}
	if matcher.matchCount >= len(matcher.participants) {
		log.Error().Int("requestedMatchCount", matcher.matchCount).Msg("Requested match count is too high!")
		return ErrMatchCountTooHigh
	}

	matcher.matchShuffle()

	// Loop through all the participants and assign matches
	// each participant's match is the next participant in the list,
	// second match is the participant after that, third, etc
	for i, p := range matcher.participants {
		giftees := []*participant.Participant{}
		for j := 1; j <= matcher.matchCount; j++ {
			giftees = append(
				giftees,
				matcher.participants[(i+j)%len(matcher.participants)],
			)
		}
		matcher.matches[p] = giftees
	}

	// TODO write matches to file

	log.Info().
		Int("participantCount", len(matcher.participants)).
		Msg("Participants shuffled, matched, and written to file")
	return nil
}

func (matcher *Matcher) matchShuffle() int {
	iterations := rand.Int31n(int32(len(matcher.participants)-1)) + 1 // At least once, but no more than the number of participants minus 1
	log.Trace().Int32("iterations", iterations).Msg("shuffled participant list")
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
