package surt

import "testing"

func TestIAURLCanonicalizer(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"http://ARCHIVE.ORG/", "http://archive.org/"},
		{"http://www.archive.org:80/", "http://archive.org/"},
		{"https://www.archive.org:80/", "https://archive.org:80/"},
		{"http://www.archive.org:443/", "http://archive.org:443/"},
		{"https://www.archive.org:443/", "https://archive.org/"},
		{"http://www.archive.org/big/", "http://archive.org/big"},
		{"dns:www.archive.org", "dns:www.archive.org"},
	}

	for _, tc := range tests {
		got, err := IACanonicalize(tc.in)
		if err != nil {
			t.Fatalf("Canonicalize(%q) returned unexpected error: %v", tc.in, err)
		}
		if got != tc.want {
			t.Errorf("Canonicalize(%q) = %q; want %q", tc.in, got, tc.want)
		}
	}
}
