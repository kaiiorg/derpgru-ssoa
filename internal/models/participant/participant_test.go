package participant

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDumpGifterMessage_AllValues(t *testing.T) {
	// Arrange
	p := &Participant{
		Timestamp:       uuid.New().String(),
		DiscordUsername: uuid.New().String(),
		Address:         uuid.New().String(),
		Name:            uuid.New().String(),
		Over21:          uuid.New().String(),
		TShirtSize:      uuid.New().String(),
		Allergies:       uuid.New().String(),
		Requests:        uuid.New().String(),
	}
	gifter := &Participant{
		DiscordUsername: uuid.New().String(),
	}
	eventName := uuid.New().String()

	// Act
	result := p.DumpGifterMessage(eventName, gifter)

	// Assert
	// We're not going to do much beyond check if the number of lines match what we expect and the final string contains all the UUIDs
	require.NotEmpty(t, result)
	require.Contains(t, result, gifter.DiscordUsername)
	require.Contains(t, result, eventName)
	require.Contains(t, result, p.DiscordUsername)
	require.Contains(t, result, p.Address)
	require.Contains(t, result, p.Name)
	require.Contains(t, result, p.Over21)
	require.Contains(t, result, p.TShirtSize)
	require.Contains(t, result, p.Allergies)
	require.Contains(t, result, p.Requests)
	require.Equal(t, 2, strings.Count(result, "```"))
	require.Equal(t, 13, strings.Count(result, "\n"))
}

func TestDumpGifterMessage_MissingName(t *testing.T) {
	// Arrange
	p := &Participant{
		Timestamp:       uuid.New().String(),
		DiscordUsername: uuid.New().String(),
		Address:         uuid.New().String(),
		Over21:          uuid.New().String(),
		TShirtSize:      uuid.New().String(),
		Allergies:       uuid.New().String(),
		Requests:        uuid.New().String(),
		//Name:          uuid.New().String(),
	}
	gifter := &Participant{
		DiscordUsername: uuid.New().String(),
	}
	eventName := uuid.New().String()

	// Act
	result := p.DumpGifterMessage(eventName, gifter)

	// Assert
	// We're not going to do much beyond check if the number of lines match what we expect and the final string contains all the UUIDs
	require.NotEmpty(t, result)
	require.Contains(t, result, gifter.DiscordUsername)
	require.Contains(t, result, eventName)
	require.Contains(t, result, p.DiscordUsername)
	require.Contains(t, result, p.Address)
	require.Contains(t, result, p.Over21)
	require.Contains(t, result, p.TShirtSize)
	require.Contains(t, result, p.Allergies)
	require.Contains(t, result, p.Requests)
	require.NotContains(t, result, "Name: ")
	require.Equal(t, 2, strings.Count(result, "```"))
	require.Equal(t, 12, strings.Count(result, "\n"))
}

func TestDumpGifterMessage_MissingAllergies(t *testing.T) {
	// Arrange
	p := &Participant{
		Timestamp:       uuid.New().String(),
		DiscordUsername: uuid.New().String(),
		Address:         uuid.New().String(),
		Name:            uuid.New().String(),
		Over21:          uuid.New().String(),
		TShirtSize:      uuid.New().String(),
		Requests:        uuid.New().String(),
		// Allergies:    uuid.New().String(),
	}
	gifter := &Participant{
		DiscordUsername: uuid.New().String(),
	}
	eventName := uuid.New().String()

	// Act
	result := p.DumpGifterMessage(eventName, gifter)

	// Assert
	// We're not going to do much beyond check if the number of lines match what we expect and the final string contains all the UUIDs
	require.NotEmpty(t, result)
	require.Contains(t, result, gifter.DiscordUsername)
	require.Contains(t, result, eventName)
	require.Contains(t, result, p.DiscordUsername)
	require.Contains(t, result, p.Address)
	require.Contains(t, result, p.Name)
	require.Contains(t, result, p.Over21)
	require.Contains(t, result, p.TShirtSize)
	require.Contains(t, result, p.Requests)
	require.NotContains(t, result, "Allergies: ")
	require.Equal(t, 2, strings.Count(result, "```"))
	require.Equal(t, 12, strings.Count(result, "\n"))
}

func TestDumpGifterMessage_Requests(t *testing.T) {
	// Arrange
	p := &Participant{
		Timestamp:       uuid.New().String(),
		DiscordUsername: uuid.New().String(),
		Address:         uuid.New().String(),
		Name:            uuid.New().String(),
		Over21:          uuid.New().String(),
		TShirtSize:      uuid.New().String(),
		Allergies:       uuid.New().String(),
		// Requests:     uuid.New().String(),
	}
	gifter := &Participant{
		DiscordUsername: uuid.New().String(),
	}
	eventName := uuid.New().String()

	// Act
	result := p.DumpGifterMessage(eventName, gifter)

	// Assert
	// We're not going to do much beyond check if the number of lines match what we expect and the final string contains all the UUIDs
	require.NotEmpty(t, result)
	require.Contains(t, result, gifter.DiscordUsername)
	require.Contains(t, result, eventName)
	require.Contains(t, result, p.DiscordUsername)
	require.Contains(t, result, p.Address)
	require.Contains(t, result, p.Name)
	require.Contains(t, result, p.Over21)
	require.Contains(t, result, p.TShirtSize)
	require.Contains(t, result, p.Allergies)
	require.NotContains(t, result, "Requests: ")
	require.Equal(t, 2, strings.Count(result, "```"))
	require.Equal(t, 11, strings.Count(result, "\n"))
}
