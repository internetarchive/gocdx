// Package surt is a best-effort reimplementation of the Sort-friendly URI Reordering Transform (SURT) python package.
//
// This essentially skips the more complex aspects of the original implementation in favor of a simpler, more idiomatic Go approach.
//
// All IA-specific cannonicalization is implemented following the tests found here https://github.com/internetarchive/surt/blob/6934c321b3e2f66af9c001d882475949f00570c5/tests/test_surt.py#L140 and the URL are parsed using the fasturl Ragel state machine.
package surt

import (
	"strings"

	"github.com/ImVexed/fasturl"
)

func Massage(url string, opts ...MassageOpts) (massagedURL string, err error) {
	parsedURL, err := fasturl.ParseURL(url)
	if err != nil {
		return "", err
	}

	// Perform any necessary massaging on the parsed URL
	var result strings.Builder

	// Example of massaging: converting to lowercase and removing trailing slashes
	result.WriteString(strings.ToLower(parsedURL.Host))
	if parsedURL.Path != "" {
		result.WriteString(parsedURL.Path)
	}
	if parsedURL.Query != "" {
		result.WriteString("?" + parsedURL.Query)
	}
	if parsedURL.Fragment != "" {
		result.WriteString("#" + parsedURL.Fragment)
	}

	return result.String(), nil
}
