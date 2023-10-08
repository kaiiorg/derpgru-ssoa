package matcher

import (
	"os"
	"testing"

	"github.com/gocarina/gocsv"
	"github.com/stretchr/testify/require"
)

func TestLoad_NoError(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.participantsFilepath = "./test_data/sanitized.csv" // 1 header + 60 entries

	// Act
	err := matcher.load(testingCmd(), []string{})

	// Assert
	require.NoError(t, err)
	require.Len(t, matcher.participants, 60)
	require.Len(t, matcher.participantsMap, 60)
}

func TestLoad_MissingParticipantsFile(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.participantsFilepath = "./test_data/nonexistent.csv"

	// Act
	err := matcher.load(testingCmd(), []string{})

	// Assert
	require.ErrorIs(t, err, os.ErrNotExist)
	require.Len(t, matcher.participants, 0)
}

func TestLoad_EmptyParticipantsFile(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.participantsFilepath = "./test_data/empty.sanitized.csv"

	// Act
	err := matcher.load(testingCmd(), []string{})

	// Assert
	require.ErrorIs(t, err, gocsv.ErrEmptyCSVFile)
	require.Len(t, matcher.participants, 0)
}

func TestLoadMatches_NoError(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.participantsFilepath = "./test_data/sanitized.csv" // 1 header + 60 entries
	err := matcher.load(testingCmd(), []string{})
	require.NoError(t, err)
	require.Len(t, matcher.participants, 60)
	matcher.matchesFilepath = "./test_data/3.matches.sanitized.csv"

	// Act
	err = matcher.loadMatches()

	// Assert
	require.NoError(t, err)
	require.Equal(t, 3, matcher.matchCount)
	require.Equal(t, len(matcher.participants), len(matcher.matches))
	for _, participant := range matcher.participants {
		// Make sure every participant in the participants file is represented in the matches map
		require.Contains(t, matcher.matches, participant.DiscordUsername)
		// Make sure each gifter has the expected number of options
		require.Len(t, matcher.matches[participant.DiscordUsername], matcher.matchCount)
		// Make sure each option for this gifter is an actual participant
		for _, option := range matcher.matches[participant.DiscordUsername] {
			require.Contains(t, matcher.participantsMap, option)
		}
	}
}

func TestLoadMatches_NoHeaderNoError(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.participantsFilepath = "./test_data/sanitized.csv" // 1 header + 60 entries
	err := matcher.load(testingCmd(), []string{})
	require.NoError(t, err)
	require.Len(t, matcher.participants, 60)
	matcher.matchesFilepath = "./test_data/noheader.3.matches.sanitized.csv"

	// Act
	err = matcher.loadMatches()

	// Assert
	require.NoError(t, err)
	require.Equal(t, 3, matcher.matchCount)
	require.Equal(t, len(matcher.participants), len(matcher.matches))
	for _, participant := range matcher.participants {
		// Make sure every participant in the participants file is represented in the matches map
		require.Contains(t, matcher.matches, participant.DiscordUsername)
		// Make sure each gifter has the expected number of options
		require.Len(t, matcher.matches[participant.DiscordUsername], matcher.matchCount)
		// Make sure each option for this gifter is an actual participant
		for _, option := range matcher.matches[participant.DiscordUsername] {
			require.Contains(t, matcher.participantsMap, option)
		}
	}
}

func TestLoadMatches_TooFewMatches(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.participantsFilepath = "./test_data/sanitized.csv" // 1 header + 2 entries
	err := matcher.load(testingCmd(), []string{})
	require.NoError(t, err)
	require.Len(t, matcher.participants, 60)
	matcher.matchesFilepath = "./test_data/ErrMatchCountTooLow.matches.sanitized.csv"

	// Act
	err = matcher.loadMatches()

	// Assert
	require.ErrorIs(t, err, ErrMatchCountTooLow)
}

func TestLoadMatches_TooManyMatches(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.participantsFilepath = "./test_data/ErrMatchCountTooHigh.sanitized.csv" // 1 header + 2 entries
	err := matcher.load(testingCmd(), []string{})
	require.NoError(t, err)
	require.Len(t, matcher.participants, 2)
	matcher.matchesFilepath = "./test_data/3.matches.sanitized.csv"

	// Act
	err = matcher.loadMatches()

	// Assert
	require.ErrorIs(t, err, ErrMatchCountTooHigh)
}

func TestLoadMatches_MissingGifter(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.participantsFilepath = "./test_data/sanitized.csv" // 1 header + 2 entries
	err := matcher.load(testingCmd(), []string{})
	require.NoError(t, err)
	require.Len(t, matcher.participants, 60)
	matcher.matchesFilepath = "./test_data/missing.gifter.matches.sanitized.csv"

	// Act
	err = matcher.loadMatches()

	// Assert
	require.ErrorIs(t, err, ErrMatchGifterNotFound)
}

func TestLoadMatches_MissingOption(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.participantsFilepath = "./test_data/sanitized.csv" // 1 header + 2 entries
	err := matcher.load(testingCmd(), []string{})
	require.NoError(t, err)
	require.Len(t, matcher.participants, 60)
	matcher.matchesFilepath = "./test_data/missing.giftee.matches.sanitized.csv"

	// Act
	err = matcher.loadMatches()

	// Assert
	require.ErrorIs(t, err, ErrMatchOptionNotFound)
}
