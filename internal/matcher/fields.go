package matcher

import (
	"bufio"
	"bytes"
	"encoding/csv"
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
	header, _, err := reader.ReadLine()
	if err != nil {
		return err
	}

	// Modify header as needed, then write it to temp file
	header, err = matcher.replaceHeaders(header)
	if err != nil {
		return err
	}

	_, err = modifiedF.Write(header)
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

func (matcher *Matcher) replaceHeaders(rawHeader []byte) ([]byte, error) {
	// Parse the header as CSV
	reader := csv.NewReader(bytes.NewBuffer(rawHeader))
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// TODO load this from file or CLI
	replacements := map[string]string{
		"Shipping Address (include country)": "Shipping Address",
		"Special requests":                   "Special requests or notes (for either your gifter or the henchmen)",
	}

	// Modify each header as needed
	resultHeaders := make([]string, len(headers))
	for i, header := range headers {
		// Check if this header is in our map of replacements
		// If it is, use the replacement value
		// If it isn't, use the current value as is
		to, found := replacements[header]
		if found {
			resultHeaders[i] = to
		} else {
			resultHeaders[i] = header
		}
	}

	// Write the header back to CSV
	buffer := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buffer)
	err = writer.Write(resultHeaders)
	if err != nil {
		return nil, err
	}
	writer.Flush()

	return buffer.Bytes(), nil
}
