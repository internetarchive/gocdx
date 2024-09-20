package gocdx

import "time"

// Record represents a single record in a CDX file
type Record struct {
	// Default fields (CDX N b a m s k r M S V g)
	MassagedURL          string    `json:"massaged_url"`           // N
	Timestamp            time.Time `json:"timestamp"`              // b
	OriginalURL          string    `json:"original_url"`           // a
	MIMEType             string    `json:"mime_type"`              // m
	StatusCode           int       `json:"status_code"`            // s
	NewStyleChecksum     string    `json:"new_style_checksum"`     // k
	Redirect             string    `json:"redirect"`               // r
	MetaTags             string    `json:"meta_tags"`              // M
	CompressedRecordSize int64     `json:"compressed_record_size"` // S
	CompressedArcOffset  int64     `json:"compressed_arc_offset"`  // V
	Filename             string    `json:"filename"`               // g

	// Optional fields
	CanonizedURL          string `json:"canonized_url,omitempty"`
	NewsGroup             string `json:"news_group,omitempty"`
	RulespaceCategory     string `json:"rulespace_category,omitempty"`
	CompressedDatOffset   int64  `json:"compressed_dat_offset,omitempty"`
	CanonizedFrame        string `json:"canonized_frame,omitempty"`
	LanguageDescription   string `json:"language_description,omitempty"`
	CanonizedHost         string `json:"canonized_host,omitempty"`
	CanonizedImage        string `json:"canonized_image,omitempty"`
	CanonizedJumpPoint    string `json:"canonized_jump_point,omitempty"`
	FBISChangedThing      string `json:"fbis_changed_thing,omitempty"`
	CanonizedLink         string `json:"canonized_link,omitempty"`
	CanonizedPath         string `json:"canonized_path,omitempty"`
	LanguageString        string `json:"language_string,omitempty"`
	CanonizedRedirect     string `json:"canonized_redirect,omitempty"`
	Uniqueness            string `json:"uniqueness,omitempty"`
	CanonizedURLOtherHref string `json:"canonized_url_other_href,omitempty"`
	CanonizedURLOtherSrc  string `json:"canonized_url_other_src,omitempty"`
	CanonizedURLScript    string `json:"canonized_url_script,omitempty"`
	OldStyleChecksum      string `json:"old_style_checksum,omitempty"`
	UncompressedDatOffset int64  `json:"uncompressed_dat_offset,omitempty"`
	IP                    string `json:"ip,omitempty"`
	Frame                 string `json:"frame,omitempty"`
	OriginalHost          string `json:"original_host,omitempty"`
	Image                 string `json:"image,omitempty"`
	OriginalJumpPoint     string `json:"original_jump_point,omitempty"`
	Link                  string `json:"link,omitempty"`
	ArcDocumentLength     int64  `json:"arc_document_length,omitempty"`
	Port                  int    `json:"port,omitempty"`
	OriginalPath          string `json:"original_path,omitempty"`
	Title                 string `json:"title,omitempty"`
	UncompressedArcOffset int64  `json:"uncompressed_arc_offset,omitempty"`
	URLOtherHref          string `json:"url_other_href,omitempty"`
	URLOtherSrc           string `json:"url_other_src,omitempty"`
	URLScript             string `json:"url_script,omitempty"`
}

// FieldIndex represents the indices of fields in the CDX file
type FieldIndex map[byte]int

// DefaultFields represents the default CDX fields in order
var DefaultFields = []byte{'N', 'b', 'a', 'm', 's', 'k', 'r', 'M', 'S', 'V', 'g'}
