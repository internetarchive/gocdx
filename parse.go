package cdx

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func Parse(r io.Reader) ([]Record, error) {
	scanner := bufio.NewScanner(r)
	var records []Record
	var fieldIndex FieldIndex

	// Check if the input is empty
	if !scanner.Scan() {
		return nil, errors.New("empty CDX file")
	}

	// Parse the header line
	header := scanner.Text()
	if header == "" {
		return nil, errors.New("empty header line in CDX file")
	}

	var err error
	fieldIndex, err = parseHeader(header)
	if err != nil {
		return nil, err // No need to wrap this error
	}

	// Read and parse the data lines
	for scanner.Scan() {
		line := scanner.Text()
		record, err := parseRecord(line, fieldIndex)
		if err != nil {
			return nil, err // No need to wrap this error
		}
		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	// Check if we have any records
	if len(records) == 0 {
		return nil, errors.New("no records found in CDX file")
	}

	return records, nil
}

func parseHeader(header string) (FieldIndex, error) {
	fields := strings.Fields(header)
	if len(fields) < 2 || fields[0] != "CDX" {
		return nil, errors.New("invalid CDX header")
	}

	fi := make(FieldIndex)
	for i, field := range fields[1:] {
		if len(field) != 1 {
			return nil, fmt.Errorf("invalid field specifier in header: %s", field)
		}
		fi[field[0]] = i
	}
	return fi, nil
}

func parseRecord(line string, fi FieldIndex) (Record, error) {
	fields := strings.Fields(line)
	if len(fields) < len(fi) {
		return Record{}, errors.New("insufficient fields in CDX record")
	}

	record := Record{}

	for key, index := range fi {
		if index >= len(fields) {
			continue
		}
		value := fields[index]
		var err error
		switch key {
		case 'N':
			record.MassagedURL = value
		case 'b':
			record.Timestamp, err = time.Parse("20060102150405", value)
			if err != nil {
				return Record{}, fmt.Errorf("invalid timestamp: %w", err)
			}
		case 'a':
			record.OriginalURL = value
		case 'm':
			record.MIMEType = value
		case 's':
			record.StatusCode, err = strconv.Atoi(value)
			if err != nil {
				return Record{}, fmt.Errorf("invalid status code: %w", err)
			}
		case 'k':
			record.NewStyleChecksum = value
		case 'r':
			record.Redirect = value
		case 'M':
			record.MetaTags = value
		case 'S':
			record.CompressedRecordSize, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return Record{}, fmt.Errorf("invalid compressed record size: %w", err)
			}
		case 'V':
			record.CompressedArcOffset, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return Record{}, fmt.Errorf("invalid compressed arc offset: %w", err)
			}
		case 'g':
			record.Filename = value
		case 'e':
			record.IP = value
		case 'h':
			record.OriginalHost = value
			// Add other optional fields here as needed
		}
	}

	return record, nil
}
