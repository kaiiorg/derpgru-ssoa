package matcher

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math/rand"
	"os"

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
	matcher.matchDoMatching()

	f, err := os.Create(matcher.matchesFilepath)
	if err != nil {
		log.Error().
			Err(err).
			Str("matchesFilepath", matcher.matchesFilepath).
			Msg("Participants shuffled and matched, but failed to open matches file")
		return err
	}
	defer f.Close()

	// gocsv doesn't support reading/writing CSVs in the manner we need for the matches, so we're using encoding/csv
	csvWriter := csv.NewWriter(f)
	defer csvWriter.Flush()
	err = matcher.matchWriteMatchesHeader(csvWriter)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Participants shuffled and matched, but failed to write header to matches file")
		return err
	}

	err = matcher.matchWriteMatches(csvWriter)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Participants shuffled and matched, but failed to write matches to matches file")
		return err
	}

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

func (matcher *Matcher) matchDoMatching() {
	// Loop through all the participants and assign matches
	// each participant's match is the next participant in the list,
	// second match is the participant after that, third, etc
	for i, p := range matcher.participants {
		giftees := []string{}
		for j := 1; j <= matcher.matchCount; j++ {
			giftees = append(
				giftees,
				matcher.participants[(i+j)%len(matcher.participants)].DiscordUsername,
			)
		}
		matcher.matches[p.DiscordUsername] = giftees
	}
}

func (matcher *Matcher) matchWriteMatchesHeader(csvWriter *csv.Writer) error {
	header := []string{"gifter"}
	for i := 1; i <= matcher.matchCount; i++ {
		header = append(header, fmt.Sprintf("option %d", i))
	}

	return csvWriter.Write(header)
}

func (matcher *Matcher) matchWriteMatches(csvWriter *csv.Writer) error {
	for gifter, giftees := range matcher.matches {
		line := []string{gifter}
		line = append(line, giftees...)
		err := csvWriter.Write(line)
		if err != nil {
			return err
		}
	}
	return nil
}
