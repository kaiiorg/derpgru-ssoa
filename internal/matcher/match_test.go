package matcher

import (
	"slices"
	"testing"

	"github.com/kaiiorg/derpgru-ssoa/internal/models/participant"

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

	// Act
	err = matcher.match(testingCmd(), []string{})

	// Assert
	require.NoError(t, err)
	for _, p := range matcher.participants {
		// Make sure the participant is assigned a list of potential giftees
		require.Containsf(t, matcher.matches, p, "map of potential giftees does not contain an entry for %s!", p.DiscordUsername)
		// Make sure the participant was not matched with themselves
		for _, match := range matcher.matches[p] {
			require.NotEqualf(t, p, match, "%s was matched with themselves", p.DiscordUsername)
		}
	}
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
	iterations := matcher.matchShuffle()

	// Assert
	require.NotZero(t, iterations)
	require.Equal(t, len(originalOrder), len(matcher.participants))
	require.NotZero(t, slices.CompareFunc(originalOrder, matcher.participants, participantCompareFunc), "participants were not shuffled")
}

// compareFunc is used for slices.CompareFunc and will return 0
// if the participants provided are the same and 1 if they are not
func participantCompareFunc(p1, p2 *participant.Participant) int {
	if p1 == p2 {
		return 0
	} else {
		return 1
	}
}
