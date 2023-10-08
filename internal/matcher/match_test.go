package matcher

import (
	"testing"
	"slices"

	"github.com/kaiiorg/derpgru-ssoa/internal/models/participant"

	"github.com/stretchr/testify/require"
)

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
	require.NotZero(t, slices.CompareFunc(originalOrder, matcher.participants, compareFunc), "participants were not shuffled")
}

// compareFunc is used for slices.CompareFunc and will return 0
// if the participants provided are the same and 1 if they are not
func compareFunc(p1, p2 *participant.Participant) int {
	if p1 == p2 {
		return 0
	} else {
		return 1
	}
}