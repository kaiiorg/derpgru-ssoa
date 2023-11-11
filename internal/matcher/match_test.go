package matcher

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/kaiiorg/derpgru-ssoa/internal/models/participant"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMatch_NoError(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.participantsFilepath = "./test_data/sanitized.csv" // 1 header + 60 entries
	err := matcher.load(testingCmd(), []string{})
	require.NoError(t, err)
	require.Len(t, matcher.participants, 60)
	matcher.matchCount = 3
	matcher.matchesFilepath = fmt.Sprintf("./test.%s.csv", uuid.New().String())
	defer os.Remove(matcher.matchesFilepath)

	// Act
	err = matcher.match(testingCmd(), []string{})

	// Assert
	require.NoError(t, err)
}

func TestMatch_MatchCountTooLow(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.participantsFilepath = "./test_data/sanitized.csv" // 1 header + 60 entries
	err := matcher.load(testingCmd(), []string{})
	require.NoError(t, err)
	require.Len(t, matcher.participants, 60)
	matcher.matchCount = 0

	// Act
	err = matcher.match(testingCmd(), []string{})

	// Assert
	require.ErrorIs(t, err, ErrMatchCountTooLow)
}

func TestMatch_MatchCountTooHigh(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.participantsFilepath = "./test_data/sanitized.csv" // 1 header + 60 entries
	err := matcher.load(testingCmd(), []string{})
	require.NoError(t, err)
	require.Len(t, matcher.participants, 60)
	matcher.matchCount = len(matcher.participants)

	// Act
	err = matcher.match(testingCmd(), []string{})

	// Assert
	require.ErrorIs(t, err, ErrMatchCountTooHigh)
}

func TestMatchShuffle(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.participantsFilepath = "./test_data/sanitized.csv" // 1 header + 60 entries
	err := matcher.load(testingCmd(), []string{})
	require.NoError(t, err)
	require.Len(t, matcher.participants, 60)
	originalOrder := slices.Clone(matcher.participants)

	// Act
	matcher.matchShuffle()

	// Assert
	require.Equal(t, len(originalOrder), len(matcher.participants))
	require.NotZero(t, slices.CompareFunc(originalOrder, matcher.participants, participantCompareFunc), "participants were not shuffled")
}

func TestMatchDoMatching(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.participantsFilepath = "./test_data/sanitized.csv" // 1 header + 60 entries
	err := matcher.load(testingCmd(), []string{})
	require.NoError(t, err)
	require.Len(t, matcher.participants, 60)
	matcher.matchCount = 3

	// Act
	matcher.matchDoMatching()

	// Assert
	for _, p := range matcher.participants {
		// Make sure the participant is assigned a list of potential giftees
		require.Containsf(t, matcher.matches, p.DiscordUsername, "map of potential giftees does not contain an entry for '%s'!", p.DiscordUsername)
		// Make sure the participant was not matched with themselves
		for _, match := range matcher.matches[p.DiscordUsername] {
			require.NotEqualf(t, p.DiscordUsername, match, "%s was matched with themselves", p.DiscordUsername)
		}
	}
}

func TestMatchWriteMatchesHeader_NoError(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.matchCount = 3
	buffer := bytes.NewBuffer([]byte{})
	csvWriter := csv.NewWriter(buffer)

	// Act
	err := matcher.matchWriteMatchesHeader(csvWriter)
	csvWriter.Flush()

	// Assert
	require.NoError(t, err)
	csvReader := csv.NewReader(buffer)
	result, err := csvReader.ReadAll()
	require.NoError(t, err)
	require.Len(t, result, 1, "only one record should have been written to the CSV file")
	require.Lenf(t, result[0], matcher.matchCount+1, "%d (matchCount + 1) fields should have been written to the record", matcher.matchCount+1)
}

func TestMatchWriteMatchesHeader_CatchesCsvError(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.matchCount = 3
	buffer := bytes.NewBuffer([]byte{})
	csvWriter := csv.NewWriter(buffer)
	csvWriter.Comma = 0 // This will cause encoding/csv to throw an error

	// Act
	err := matcher.matchWriteMatchesHeader(csvWriter)

	// Assert
	require.ErrorContains(t, err, "csv: invalid field or comment delimiter")
}

func TestMatchWriteMatches_NoError(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.matchCount = 3
	matcher.participantsFilepath = "./test_data/sanitized.csv" // 1 header + 60 entries
	err := matcher.load(testingCmd(), []string{})
	require.NoError(t, err)
	require.Len(t, matcher.participants, 60)
	matcher.matchDoMatching()
	buffer := bytes.NewBuffer([]byte{})
	csvWriter := csv.NewWriter(buffer)

	// Act
	err = matcher.matchWriteMatches(csvWriter)
	csvWriter.Flush()

	// Assert
	require.NoError(t, err)
	csvReader := csv.NewReader(buffer)
	result, err := csvReader.ReadAll()
	require.NoError(t, err)
	require.Lenf(t, result, len(matcher.participants), "%d records should have been written to the CSV file", len(matcher.participants))
	for i, record := range result {
		require.Lenf(t, record, matcher.matchCount+1, "%d (matchCount + 1) fields should have been written to the record on line %d", matcher.matchCount+1, i+1)
	}
}

func TestMatchWriteMatches_CatchesCsvError(t *testing.T) {
	// Arrange
	matcher := New()
	matcher.matchCount = 3
	matcher.participantsFilepath = "./test_data/sanitized.csv" // 1 header + 60 entries
	err := matcher.load(testingCmd(), []string{})
	require.NoError(t, err)
	require.Len(t, matcher.participants, 60)
	matcher.matchDoMatching()
	buffer := bytes.NewBuffer([]byte{})
	csvWriter := csv.NewWriter(buffer)
	csvWriter.Comma = 0 // This will cause encoding/csv to throw an error

	// Act
	err = matcher.matchWriteMatches(csvWriter)

	// Assert
	require.ErrorContains(t, err, "csv: invalid field or comment delimiter")
}

// participantCompareFunc is used for slices.CompareFunc and will return 0
// if the participants provided are the same and 1 if they are not
func participantCompareFunc(p1, p2 *participant.Participant) int {
	if p1 == p2 {
		return 0
	} else {
		return 1
	}
}
