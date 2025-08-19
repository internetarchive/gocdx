package gocdx

import (
	"bufio"
	"crypto/sha1"
	"encoding/base32"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/internetarchive/gocdx/pkg/surt"
	warc "github.com/internetarchive/gowarc"
)

type overloadedWARCRecord struct {
	*warc.Record

	compByteOffset int64
	compByteLength int64
	httpMessage    string
	httpHeaders    map[string]string
	warcFileName   string
}

// Generate reads a WARC file from the provided reader and returns a slice of Record generated from the given WARC records.
func Generate(warcFile io.ReadCloser, header string) ([]*Record, error) {
	warcReader, err := warc.NewReader(warcFile)
	if err != nil {
		return nil, err
	}

	var currentPosition int64
	var warcFileName string
	var i int

	var warcRecords []*overloadedWARCRecord
	for {
		warcRecord, size, err := warcReader.ReadRecord()
		if err != nil {
			return nil, err
		}
		if size == 0 {
			// EOF reached, no more records
			break
		}

		if i == 0 && warcRecord.Header.Get("WARC-Filename") != "" {
			warcFileName = warcRecord.Header.Get("WARC-Filename")
		} else if i == 0 {
			// If the first record does not have a WARC-Filename, this is an error that we want to soft-fail
			warcFileName = "unknown.warc"
		}

		var httpMessage string
		var httpHeaders map[string]string
		if warcRecord.Header.Get("WARC-Type") == "response" {
			httpMessage, httpHeaders = parseHTTPHeadersFromWARCRecord(warcRecord)
		}
		parsedWARCRecord := &overloadedWARCRecord{
			Record:         warcRecord,
			compByteOffset: currentPosition,
			compByteLength: size,
			httpMessage:    httpMessage,
			httpHeaders:    httpHeaders,
			warcFileName:   warcFileName,
		}

		// spew.Dump(parsedWARCRecord)

		warcRecords = append(warcRecords, parsedWARCRecord)
		parsedWARCRecord.Record.Content.Close() // Close the content to avoid memory leaks as we are not using it here

		currentPosition += size
	}

	// spew.Dump(warcRecords)

	headerFields := strings.FieldsSeq(header)

	records := make([]*Record, 0)
	for _, warcRecord := range warcRecords {
		if warcRecord.Record.Header.Get("WARC-Type") != "response" {
			// Only process records with WARC-Type "response"
			continue
		}

		record := &Record{}
		for field := range headerFields {
			switch field {
			case "N":
				record.MassagedURL = surt.Massage(warcRecord.Record.Header.Get("WARC-Target-URI"))
			case "b":
				parsedTime, err := time.Parse(time.RFC3339, warcRecord.Record.Header.Get("WARC-Date"))
				if err != nil {
					parsedTime = time.Time{}
				}
				record.Timestamp = parsedTime
			case "a":
				record.OriginalURL = warcRecord.Record.Header.Get("WARC-Target-URI")
			case "m":
				record.MIMEType = strings.TrimSuffix(warcRecord.Record.Header.Get("Content-Type"), "; msgtype=response")
			case "s":
				record.StatusCode = -1

				splittedHTTPMessage := strings.Split(warcRecord.httpMessage, " ")
				if len(splittedHTTPMessage) >= 2 {
					parsedStatusCode, err := strconv.Atoi(splittedHTTPMessage[1])
					if err == nil {
						record.StatusCode = parsedStatusCode
					}
				}
			case "k":
				trimmed := strings.TrimPrefix(warcRecord.Record.Header.Get("WARC-Block-Digest"), "sha1:")
				if trimmed != warcRecord.Record.Header.Get("WARC-Block-Digest") {
					record.NewStyleChecksum = trimmed
				} else {
					hasher := sha1.New()
					warcRecord.Record.Content.Seek(0, 0)
					io.Copy(hasher, warcRecord.Record.Content)
					record.NewStyleChecksum = base32.StdEncoding.EncodeToString(hasher.Sum(nil))
				}
			case "r":
				// TODO : clarify with whoever what to do with this field
				record.Redirect = "-"
			case "M":
				// TODO : let's ignore this field for now
				record.MetaTags = "-"
			case "S":
				record.CompressedRecordSize = warcRecord.compByteLength
			case "V":
				record.CompressedArcOffset = warcRecord.compByteOffset
			case "g":
				record.Filename = warcRecord.warcFileName
			}
		}
		records = append(records, record)
	}

	return records, nil
}

// FormatCDX formats a Record into a CDX string based on the header format.
func (r *Record) FormatCDX(header string) (string, error) {
	var result strings.Builder
	headerFields := strings.FieldsSeq(header)

	for field := range headerFields {
		switch field {
		case "N":
			result.WriteString(r.MassagedURL)
		case "b":
			result.WriteString(r.Timestamp.Format("20060102150405"))
		case "a":
			result.WriteString(r.OriginalURL)
		case "m":
			result.WriteString(r.MIMEType)
		case "s":
			result.WriteString(strconv.Itoa(r.StatusCode))
		case "k":
			result.WriteString(r.NewStyleChecksum)
		case "r":
			result.WriteString(r.Redirect)
		case "M":
			result.WriteString(r.MetaTags)
		case "S":
			result.WriteString(strconv.FormatInt(r.CompressedRecordSize, 10))
		case "V":
			result.WriteString(strconv.FormatInt(r.CompressedArcOffset, 10))
		case "g":
			result.WriteString(r.Filename)
		}
		result.WriteString(" ")
	}

	return strings.TrimSpace(result.String()), nil
}

func parseHTTPHeadersFromWARCRecord(warcRecord *warc.Record) (message string, headers map[string]string) {
	headers = make(map[string]string)

	var i int
	scanner := bufio.NewScanner(warcRecord.Content)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		// handle HTTP message
		if i == 0 {
			message = strings.Clone(line)
			i++
			continue
		}

		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			headers[strings.ToLower(parts[0])] = strings.ToLower(parts[1])
		}

		i++
	}

	return
}
