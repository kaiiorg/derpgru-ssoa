package matcher

import (
	"testing"
	"os"

	"github.com/stretchr/testify/require"
	"github.com/gocarina/gocsv"
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
