package matcher

import (
	"encoding/csv"
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

	// Copy the participants into a map so we can easily find them by username. We need them
	// in the array as well for easy index based selection for the match command to shuffle and select
	for _, participant := range matcher.participants {
		matcher.participantsMap[participant.DiscordUsername] = participant
	}

	log.Info().Str("participantsFilepath", matcher.participantsFilepath).Msg("loaded participants file")

	// If this is the generate command, load the match file too
	if cmd.Use == GENERATE_CMD_NAME {
		return matcher.loadMatches()
	}

	return nil
}

func (matcher *Matcher) loadMatches() error {
	matchesF, err := os.Open(matcher.matchesFilepath)
	if err != nil {
		return err
	}
	defer matchesF.Close()

	csvReader := csv.NewReader(matchesF)
	allMatches, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	for i, record := range allMatches {
		if i == 0 {
			// Min = gifter column + 1 giftee option
			if len(record) < 2 {
				log.Error().Int("count", len(record)).Strs("record", record).Msg("Too few fields in record")
				return ErrMatchCountTooLow
			}
			// Max = gifter column + all the other giftee options
			if len(record) > len(matcher.participants)+1 {
				log.Error().Int("count", len(record)).Strs("record", record).Msg("Too many fields in record for number of participants")
				return ErrMatchCountTooHigh
			}
			// The first record will tell us how many matches were intially configured
			matcher.matchCount = len(record) - 1 // subtract 1 to account for the gifter column

			// If the very first record isn't MATCH_FILE_HEADER_GIFTER, we assume the header was
			// left out after the event coordinators made manual adjustments. This will break
			// if there happens to be a participant who's username is MATCH_FILE_HEADER_GIFTER,
			// but we'll worry about that when such a user joins our discord server
			if record[0] == MATCH_FILE_HEADER_GIFTER {
				continue
			}
		}

		// Make sure the gifter and all of their options are all valid participants
		if _, found := matcher.participantsMap[record[0]]; !found {
			log.Error().Str("gifter", record[0]).Strs("options", record[1:]).Msg("Did not find gifter in participants file")
			return ErrMatchGifterNotFound
		}
		for _, option := range record[1:] {
			if _, found := matcher.participantsMap[option]; !found {
				log.Error().Str("gifter", record[0]).Str("option", option).Strs("options", record[1:]).Msg("Did not find option for gifter in participants file")
				return ErrMatchOptionNotFound
			}
		}

		matcher.matches[record[0]] = record[1:]
	}
	return nil
}
