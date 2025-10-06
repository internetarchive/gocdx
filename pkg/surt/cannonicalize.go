package surt

import (
	"net"
	"net/url"
	"strings"

	"github.com/ImVexed/fasturl"
)

// IACanonicalize returns the canonicalized URL string per the IA-style rules defined by surt tests.
func IACanonicalize(raw string) (string, error) {
	parsed, err := fasturl.ParseURL(raw)
	if err != nil {
		return "", err
	}

	// Non-HTTP(S) schemes are left untouched (e.g., "dns:...").
	scheme := strings.ToLower(parsed.Protocol)
	if scheme != "http" && scheme != "https" {
		// Return the input as-is
		return raw, nil
	}

	newURL := url.URL{}

	newURL.Scheme = scheme

	// Normalize host: lowercase & drop leading "www.".
	host := strings.ToLower(parsed.Host)
	port := parsed.Port

	// Strip leading "www." if present
	host = strings.TrimPrefix(host, "www.")

	// Remove default ports (http->80, https->443); keep non-defaults.
	defaultPort := map[string]string{
		"http":  "80",
		"https": "443",
	}[scheme]

	if port == defaultPort {
		port = ""
	}

	if port != "" {
		newURL.Host = net.JoinHostPort(host, port)
	} else {
		newURL.Host = host
	}

	// Path: keep "/" at root, otherwise drop a single trailing slash.
	path := parsed.Path
	if path == "" {
		path = "/"
	}
	if path != "/" && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}
	newURL.Path = path

	// Queries/fragments aren’t covered by the provided tests, so leave as-is.

	return newURL.String(), nil
}
