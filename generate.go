package gocdx

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/base32"
	"io"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/internetarchive/gocdx/pkg/surt"
	warc "github.com/internetarchive/gowarc"
	"golang.org/x/sync/errgroup"
)

type overloadedWARCRecord struct {
	*warc.Record

	compByteOffset int64
	compByteLength int64
	warcFileName   string
}

// contains helper
func hasField(fields []string, want string) bool {
	for _, f := range fields {
		if f == want {
			return true
		}
	}
	return false
}

// Generate reads a WARC file from the provided reader and returns a slice of Record
// generated from the given WARC response records. Reading from the WARC reader is
// strictly single-threaded; record processing is concurrent.
func Generate(warcFile io.Reader, header string) ([]*Record, error) {
	warcReader, err := warc.NewReader(warcFile)
	if err != nil {
		return nil, err
	}
	defer warcReader.Close()

	// Which fields are requested?
	headerFields := strings.Fields(header)
	needStatus := hasField(headerFields, "s") // requires parsing HTTP status line
	needDigest := hasField(headerFields, "k") // may require reading full record content

	// ---- Stage 1: single-threaded collection (the only critical section) ----
	var (
		warcFileName string
		warcRecords  []*overloadedWARCRecord // only "response" records are collected
	)

	for {
		rec, err := warcReader.ReadRecord()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Best effort: close any content we already collected
			for _, r := range warcRecords {
				if r.Record != nil && r.Record.Content != nil {
					_ = r.Record.Content.Close()
				}
			}
			return nil, err
		}

		switch rec.Header.Get("WARC-Type") {
		case "warcinfo":
			warcFileName = rec.Header.Get("WARC-Filename")
		case "response":
			// Keep Content open; workers will parse and then close it later.
			warcRecords = append(warcRecords, &overloadedWARCRecord{
				Record:         rec,
				compByteOffset: rec.Offset,
				compByteLength: rec.Size,
				warcFileName:   warcFileName,
			})
		default:
		}
	}

	// Nothing to do?
	if len(warcRecords) == 0 {
		return []*Record{}, nil
	}

	// ---- Stage 2: concurrent processing of collected response records ----
	records := make([]*Record, len(warcRecords))

	workers := runtime.GOMAXPROCS(0) // reasonable default; tune if needed
	g, _ := errgroup.WithContext(context.Background())

	// Work distribution: each worker pulls indexes from a channel
	type job struct{ idx int }
	jobs := make(chan job)

	// Start workers
	for w := 0; w < workers; w++ {
		g.Go(func() error {
			for j := range jobs {
				wr := warcRecords[j.idx]
				rec := &Record{}

				// Build the record according to requested fields only.
				for _, field := range headerFields {
					switch field {
					case "N":
						rec.MassagedURL = surt.Massage(wr.Record.Header.Get("WARC-Target-URI"))

					case "b":
						parsedTime, err := time.Parse(time.RFC3339, wr.Record.Header.Get("WARC-Date"))
						if err != nil {
							parsedTime = time.Time{}
						}
						rec.Timestamp = parsedTime

					case "a":
						rec.OriginalURL = wr.Record.Header.Get("WARC-Target-URI")

					case "m":
						rec.MIMEType = strings.TrimSuffix(
							wr.Record.Header.Get("Content-Type"),
							"; msgtype=response",
						)

					case "s":
						// Only parse HTTP status if requested
						rec.StatusCode = -1
						if needStatus {
							httpMessage, _ := parseHTTPHeadersFromWARCRecord(wr.Record)
							parts := strings.Split(httpMessage, " ")
							if len(parts) >= 2 {
								if sc, err := strconv.Atoi(parts[1]); err == nil {
									rec.StatusCode = sc
								}
							}
							// Reset content for any subsequent reads (e.g., digest)
							if rs, ok := wr.Record.Content.(io.ReadSeeker); ok {
								_, _ = rs.Seek(0, io.SeekStart)
							}
						}

					case "k":
						// Prefer the digest from header if present; otherwise compute it.
						trimmed := strings.TrimPrefix(wr.Record.Header.Get("WARC-Block-Digest"), "sha1:")
						if trimmed != wr.Record.Header.Get("WARC-Block-Digest") {
							rec.NewStyleChecksum = trimmed
						} else if needDigest {
							hasher := sha1.New()
							// Ensure we start from the beginning (status parsing may have read some bytes)
							if rs, ok := wr.Record.Content.(io.ReadSeeker); ok {
								_, _ = rs.Seek(0, io.SeekStart)
							}
							if _, err := io.Copy(hasher, wr.Record.Content); err != nil {
								// Close before returning error
								_ = wr.Record.Content.Close()
								return err
							}
							rec.NewStyleChecksum = base32.StdEncoding.EncodeToString(hasher.Sum(nil))
							// Reset again is not necessary since we'll close Content below
						}

					case "r":
						// TODO: clarify; keep placeholder
						rec.Redirect = "-"

					case "M":
						// TODO: ignore for now
						rec.MetaTags = "-"

					case "S":
						rec.CompressedRecordSize = wr.compByteLength

					case "V":
						rec.CompressedArcOffset = wr.compByteOffset

					case "g":
						rec.Filename = wr.warcFileName
					}
				}

				// Release resources for this record
				if wr.Record != nil && wr.Record.Content != nil {
					_ = wr.Record.Content.Close()
				}

				records[j.idx] = rec
			}
			return nil
		})
	}

	// Feed jobs
	go func() {
		for i := range warcRecords {
			jobs <- job{idx: i}
		}
		close(jobs)
	}()

	// Wait for all workers
	if err := g.Wait(); err != nil {
		return nil, err
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
		case "CDX":
			continue
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
