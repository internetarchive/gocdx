package gocdx

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Record
		wantErr  bool
		errMsg   string
	}{
		{
			name: "Parse default fields",
			input: `CDX N b a m s k r M S V g
com,example)/ 20230731193152 https://example.com/ text/html 200 K5UZWMOAOHAFAVNO2QBMBBCMGAAC7K3J - - 1234 5678 example.warc.gz
`,
			expected: []Record{
				{
					MassagedURL:          "com,example)/",
					Timestamp:            time.Date(2023, 7, 31, 19, 31, 52, 0, time.UTC),
					OriginalURL:          "https://example.com/",
					MIMEType:             "text/html",
					StatusCode:           200,
					NewStyleChecksum:     "K5UZWMOAOHAFAVNO2QBMBBCMGAAC7K3J",
					Redirect:             "-",
					MetaTags:             "-",
					CompressedRecordSize: 1234,
					CompressedArcOffset:  5678,
					Filename:             "example.warc.gz",
				},
			},
			wantErr: false,
		},
		{
			name: "Parse with optional fields",
			input: `CDX N b a m s k r M S V g e h
com,example)/ 20230731193152 https://example.com/ text/html 200 K5UZWMOAOHAFAVNO2QBMBBCMGAAC7K3J - - 1234 5678 example.warc.gz 192.168.1.1 example.com
`,
			expected: []Record{
				{
					MassagedURL:          "com,example)/",
					Timestamp:            time.Date(2023, 7, 31, 19, 31, 52, 0, time.UTC),
					OriginalURL:          "https://example.com/",
					MIMEType:             "text/html",
					StatusCode:           200,
					NewStyleChecksum:     "K5UZWMOAOHAFAVNO2QBMBBCMGAAC7K3J",
					Redirect:             "-",
					MetaTags:             "-",
					CompressedRecordSize: 1234,
					CompressedArcOffset:  5678,
					Filename:             "example.warc.gz",
					IP:                   "192.168.1.1",
					OriginalHost:         "example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "Parse multiple records",
			input: `CDX N b a m s k r M S V g
com,example)/ 20230731193152 https://example.com/ text/html 200 K5UZWMOAOHAFAVNO2QBMBBCMGAAC7K3J - - 1234 5678 example1.warc.gz
org,example)/ 20230731193153 https://example.org/ text/plain 404 L6VZXNPAPIRGBWOP3RCNCDNDHABD8L4K - - 5678 9012 example2.warc.gz
`,
			expected: []Record{
				{
					MassagedURL:          "com,example)/",
					Timestamp:            time.Date(2023, 7, 31, 19, 31, 52, 0, time.UTC),
					OriginalURL:          "https://example.com/",
					MIMEType:             "text/html",
					StatusCode:           200,
					NewStyleChecksum:     "K5UZWMOAOHAFAVNO2QBMBBCMGAAC7K3J",
					Redirect:             "-",
					MetaTags:             "-",
					CompressedRecordSize: 1234,
					CompressedArcOffset:  5678,
					Filename:             "example1.warc.gz",
				},
				{
					MassagedURL:          "org,example)/",
					Timestamp:            time.Date(2023, 7, 31, 19, 31, 53, 0, time.UTC),
					OriginalURL:          "https://example.org/",
					MIMEType:             "text/plain",
					StatusCode:           404,
					NewStyleChecksum:     "L6VZXNPAPIRGBWOP3RCNCDNDHABD8L4K",
					Redirect:             "-",
					MetaTags:             "-",
					CompressedRecordSize: 5678,
					CompressedArcOffset:  9012,
					Filename:             "example2.warc.gz",
				},
			},
			wantErr: false,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: nil,
			wantErr:  true,
			errMsg:   "empty CDX file",
		},
		{
			name:     "Only header, no records",
			input:    "CDX N b a m s k r M S V g\n",
			expected: nil,
			wantErr:  true,
			errMsg:   "no records found in CDX file",
		},
		{
			name:     "Invalid header",
			input:    "INVALID HEADER\n",
			expected: nil,
			wantErr:  true,
			errMsg:   "invalid CDX header",
		},
		{
			name: "Invalid record (not enough fields)",
			input: `CDX N b a m s k r M S V g
com,example)/ 20230731193152
`,
			expected: nil,
			wantErr:  true,
			errMsg:   "insufficient fields in CDX record",
		},
		{
			name: "Invalid timestamp",
			input: `CDX N b a m s k r M S V g
com,example)/ INVALIDTIME https://example.com/ text/html 200 K5UZWMOAOHAFAVNO2QBMBBCMGAAC7K3J - - 1234 5678 example.warc.gz
`,
			expected: nil,
			wantErr:  true,
			errMsg:   "invalid timestamp: parsing time \"INVALIDTIME\" as \"20060102150405\": cannot parse",
		},
		{
			name: "Invalid status code",
			input: `CDX N b a m s k r M S V g
com,example)/ 20230731193152 https://example.com/ text/html INVALID K5UZWMOAOHAFAVNO2QBMBBCMGAAC7K3J - - 1234 5678 example.warc.gz
`,
			expected: nil,
			wantErr:  true,
			errMsg:   "invalid status code",
		},
		{
			name: "Invalid compressed record size",
			input: `CDX N b a m s k r M S V g
com,example)/ 20230731193152 https://example.com/ text/html 200 K5UZWMOAOHAFAVNO2QBMBBCMGAAC7K3J - - INVALID 5678 example.warc.gz
`,
			expected: nil,
			wantErr:  true,
			errMsg:   "invalid compressed record size",
		},
		{
			name: "Invalid compressed arc offset",
			input: `CDX N b a m s k r M S V g
com,example)/ 20230731193152 https://example.com/ text/html 200 K5UZWMOAOHAFAVNO2QBMBBCMGAAC7K3J - - 1234 INVALID example.warc.gz
`,
			expected: nil,
			wantErr:  true,
			errMsg:   "invalid compressed arc offset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewBufferString(tt.input)
			got, err := Parse(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Parse() error message = %v, want it to contain %v", err.Error(), tt.errMsg)
				}
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Parse() got = %v, want %v", got, tt.expected)
			}
		})
	}
}
