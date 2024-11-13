package matcher

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/kaiiorg/derpgru-ssoa/internal/models/participant"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func (matcher *Matcher) fields(cmd *cobra.Command, args []string) {
	fields := reflect.VisibleFields(reflect.TypeOf(participant.Participant{}))
	for _, field := range fields {
		csvTag, found := field.Tag.Lookup("csv")
		if !found {
			continue
		}

		log.Info().Str("csv", csvTag).Str("field", field.Name).Send()
	}
}

func (matcher *Matcher) modify(cmd *cobra.Command, args []string) error {
	// Open new file next to CSV file
	dir, file := filepath.Split(matcher.participantsFilepath)
	file = fmt.Sprintf("modified.%s", file)
	modifiedFilepath := filepath.Join(dir, file)

	modifiedF, err := os.Create(modifiedFilepath)
	if err != nil {
		return err
	}
	defer modifiedF.Close()

	// Open CSV file
	participantsF, err := os.Open(matcher.participantsFilepath)
	if err != nil {
		return err
	}
	defer participantsF.Close()
	reader := bufio.NewReader(participantsF)

	// Read first line to get the header
	headerBytes, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	header := string(headerBytes)

	// Modify header as needed, save to temp file

	// TODO actually modify it
	header += "\n"

	_, err = modifiedF.WriteString(header)
	if err != nil {
		return err
	}

	// Read remaining file to temp file
	_, err = io.Copy(modifiedF, reader)
	if err != nil {
		return err
	}

	log.Info().Msg("Modify works")
	return nil
}
