package matcher

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/kaiiorg/derpgru-ssoa/internal/models/participant"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func (matcher *Matcher) reflectFields() map[string]string {
	f := map[string]string{}
	fields := reflect.VisibleFields(reflect.TypeOf(participant.Participant{}))
	for _, field := range fields {
		csvTag, found := field.Tag.Lookup("csv")
		if !found {
			continue
		}
		f[field.Name] = csvTag
	}
	return f
}

func (matcher *Matcher) fields(cmd *cobra.Command, args []string) {
	fields := matcher.reflectFields()
	for fieldName, fieldCsv := range fields {
		log.Info().Str("csv", fieldCsv).Str("field", fieldName).Send()
	}
}

func (matcher *Matcher) modify(cmd *cobra.Command, args []string) error {
	key, err := matcher.loadModifyKeyFromFile()
	if err != nil {
		return err
	}

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
	header, err = matcher.replaceHeaders(key, header)
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

	return nil
}

func (matcher *Matcher) loadModifyKeyFromFile() (map[string]string, error) {
	// Attempt to read from the key file
	keyJson, err := os.ReadFile(matcher.keyModifyKeyFilepath)
	if err != nil {
		// If the error was something other than it doesn't exist, return the error now
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		// Write an example file, then return an error
		fields := map[string]string{}
		for _, csvField := range matcher.reflectFields() {
			fields[csvField] = csvField
		}
		keyJson, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return nil, errors.Join(
				errors.New("key file not found"),
				errors.New("failed to marshal example json"),
				err,
			)
		}
		err = os.WriteFile(matcher.keyModifyKeyFilepath, keyJson, os.ModePerm)
		if err != nil {
			return nil, errors.Join(
				errors.New("key file not found"),
				errors.New("failed to write example file"),
				err,
			)
		}
		return nil, fmt.Errorf("key file not found, wrote example to %s", matcher.keyModifyKeyFilepath)
	}

	key := map[string]string{}
	err = json.Unmarshal(keyJson, &key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (matcher *Matcher) replaceHeaders(key map[string]string, rawHeader []byte) ([]byte, error) {
	// Parse the header as CSV
	reader := csv.NewReader(bytes.NewBuffer(rawHeader))
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Modify each header as needed
	resultHeaders := make([]string, len(headers))
	for i, header := range headers {
		// Check if this header is in our map of replacements
		// If it is, use the replacement value
		// If it isn't, use the current value as is
		to, found := key[header]
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
