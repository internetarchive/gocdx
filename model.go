package gocdx

import "time"

// Record represents a single record in a CDX file.
// Default fields set by the Internet Archive are (CDX N b a m s k r M S V g) as per IIPC specifications.
type Record struct {
	CanonizedURL                   string    `json:"canonized_url,omitempty" cdx:"A"`            // A
	NewsGroup                      string    `json:"news_group,omitempty" cdx:"B"`               // B
	RulespaceCategory              string    `json:"rulespace_category,omitempty" cdx:"C"`       // C
	CompressedDatOffset            int64     `json:"compressed_dat_offset,omitempty" cdx:"D"`    // D
	CanonizedFrame                 string    `json:"canonized_frame,omitempty" cdx:"F"`          // F
	MultiColumnLanguageDescription string    `json:"language_description,omitempty" cdx:"G"`     // G
	CanonizedHost                  string    `json:"canonized_host,omitempty" cdx:"H"`           // H
	CanonizedImage                 string    `json:"canonized_image,omitempty" cdx:"I"`          // I
	CanonizedJumpPoint             string    `json:"canonized_jump_point,omitempty" cdx:"J"`     // J
	FBISChangedThing               string    `json:"fbis_changed_thing,omitempty" cdx:"K"`       // K
	CanonizedLink                  string    `json:"canonized_link,omitempty" cdx:"L"`           // L
	MetaTags                       string    `json:"meta_tags" cdx:"M"`                          // M
	MassagedURL                    string    `json:"massaged_url" cdx:"N"`                       // N
	CanonizedPath                  string    `json:"canonized_path,omitempty" cdx:"P"`           // P
	LanguageString                 string    `json:"language_string,omitempty" cdx:"Q"`          // Q
	CanonizedRedirect              string    `json:"canonized_redirect,omitempty" cdx:"R"`       // R
	CompressedRecordSize           int64     `json:"compressed_record_size" cdx:"S"`             // S
	Uniqueness                     string    `json:"uniqueness,omitempty" cdx:"U"`               // U
	CompressedArcOffset            int64     `json:"compressed_arc_offset" cdx:"V"`              // V
	CanonizedURLOtherHref          string    `json:"canonized_url_other_href,omitempty" cdx:"X"` // X
	CanonizedURLOtherSrc           string    `json:"canonized_url_other_src,omitempty" cdx:"Y"`  // Y
	CanonizedURLScript             string    `json:"canonized_url_script,omitempty" cdx:"Z"`     // Z
	OriginalURL                    string    `json:"original_url" cdx:"a"`                       // a
	Timestamp                      time.Time `json:"timestamp" cdx:"b"`                          // b
	OldStyleChecksum               string    `json:"old_style_checksum,omitempty" cdx:"c"`       // c
	UncompressedDatOffset          int64     `json:"uncompressed_dat_offset,omitempty" cdx:"d"`  // d
	IP                             string    `json:"ip,omitempty" cdx:"e"`                       // e
	Frame                          string    `json:"frame,omitempty" cdx:"f"`                    // f
	Filename                       string    `json:"filename" cdx:"g"`                           // g
	OriginalHost                   string    `json:"original_host,omitempty" cdx:"h"`            // h
	Image                          string    `json:"image,omitempty" cdx:"i"`                    // i
	OriginalJumpPoint              string    `json:"original_jump_point,omitempty" cdx:"j"`      // j
	NewStyleChecksum               string    `json:"new_style_checksum" cdx:"k"`                 // k
	Link                           string    `json:"link,omitempty" cdx:"l"`                     // l
	MIMEType                       string    `json:"mime_type" cdx:"m"`                          // m
	ArcDocumentLength              int64     `json:"arc_document_length,omitempty" cdx:"n"`      // n
	Port                           int       `json:"port,omitempty" cdx:"o"`                     // o
	OriginalPath                   string    `json:"original_path,omitempty" cdx:"p"`            // p
	Redirect                       string    `json:"redirect" cdx:"r"`                           // r
	StatusCode                     int       `json:"status_code" cdx:"s"`                        // s
	Title                          string    `json:"title,omitempty" cdx:"t"`                    // t
	UncompressedArcOffset          int64     `json:"uncompressed_arc_offset,omitempty" cdx:"v"`  // v
	URLOtherHref                   string    `json:"url_other_href,omitempty" cdx:"x"`           // x
	URLOtherSrc                    string    `json:"url_other_src,omitempty" cdx:"y"`            // y
	URLScript                      string    `json:"url_script,omitempty" cdx:"z"`               // z
}

// FieldIndex represents the indices of fields in the CDX file
type FieldIndex map[byte]int

// DefaultFields represents the defaultorder
var DefaultFields = []byte{'N', 'b', 'a', 'm', 's', 'k', 'r', 'M', 'S', 'V', 'g'}
