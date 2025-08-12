package surt

import "testing"

func TestSURT(t *testing.T) {
	tests := []struct {
		in   string
		opts []MassageOpts
		want string
	}{
		{in: "", want: "-"},
		{in: "filedesc:foo.arc.gz", want: "filedesc:foo.arc.gz"},
		{in: "filedesc:/foo.arc.gz", want: "filedesc:/foo.arc.gz"},
		{in: "filedesc://foo.arc.gz", want: "filedesc://foo.arc.gz"},
		{in: "warcinfo:foo.warc.gz", want: "warcinfo:foo.warc.gz"},
		{in: "dns:alexa.com", want: "dns:alexa.com"},
		{in: "dns:archive.org", want: "dns:archive.org"},

		{in: "http://www.archive.org/", want: "org,archive)/"},
		{in: "http://archive.org/", want: "org,archive)/"},
		{in: "http://archive.org/goo/", want: "org,archive)/goo"},
		{in: "http://archive.org/goo/?", want: "org,archive)/goo"},
		{in: "http://archive.org/goo/?b&a", want: "org,archive)/goo?a&b"},
		{in: "http://archive.org/goo/?a=2&b&a=1", want: "org,archive)/goo?a=1&a=2&b"},

		// trailing comma mode
		{in: "http://archive.org/goo/?a=2&b&a=1", opts: []MassageOpts{WithTrailingComma}, want: "org,archive,)/goo?a=1&a=2&b"},
		{in: "dns:archive.org", opts: []MassageOpts{WithTrailingComma}, want: "dns:archive.org"},
		{in: "warcinfo:foo.warc.gz", opts: []MassageOpts{WithTrailingComma}, want: "warcinfo:foo.warc.gz"},

		// PHP session id:
		{in: "http://archive.org/index.php?PHPSESSID=0123456789abcdefghijklemopqrstuv&action=profile;u=4221", want: "org,archive)/index.php?action=profile;u=4221"},

		// WHOIS url:
		{in: "whois://whois.isoc.org.il/shaveh.co.il", want: "il,org,isoc,whois)/shaveh.co.il"},

		// Simple customization
		{in: "http://www.example.com/", opts: []MassageOpts{WithScheme}, want: "http://(com,example)/"},
		{in: "http://www.example.com/", opts: []MassageOpts{}, want: "com,example)/"},
		{in: "http://www.example.com/", opts: []MassageOpts{WithScheme, WithTrailingComma}, want: "http://(com,example,)/"},
		{in: "https://www.example.com/", opts: []MassageOpts{WithScheme, WithTrailingComma}, want: "https://(com,example,)/"},
		{in: "ftp://www.example.com/", opts: []MassageOpts{WithTrailingComma}, want: "com,example,)/"},
		{in: "ftp://www.example.com/", opts: []MassageOpts{}, want: "com,example)/"},
		{in: "ftp://www.example.com/", opts: []MassageOpts{WithScheme, WithTrailingComma}, want: "ftp://(com,example,)/"},
		{in: "http://www.example.com/", opts: []MassageOpts{WithScheme, NoHostMassage}, want: "http://(com,example,www)/"},
		{in: "http://www.example.com/", opts: []MassageOpts{NoHostMassage}, want: "com,example,www)/"},
		{in: "http://www.example.com/", opts: []MassageOpts{WithScheme, WithTrailingComma, NoHostMassage}, want: "http://(com,example,www,)/"},
		{in: "https://www.example.com/", opts: []MassageOpts{WithScheme, WithTrailingComma, NoHostMassage}, want: "https://(com,example,www,)/"},
		{in: "ftp://www.example.com/", opts: []MassageOpts{WithScheme, WithTrailingComma, NoHostMassage}, want: "ftp://(com,example,www,)/"},

		{in: "mailto:foo@example.com", opts: []MassageOpts{WithScheme}, want: "mailto:foo@example.com"},
		{in: "mailto:foo@example.com", opts: []MassageOpts{WithTrailingComma}, want: "mailto:foo@example.com"},
		{in: "mailto:foo@example.com", opts: []MassageOpts{WithScheme, WithTrailingComma}, want: "mailto:foo@example.com"},
		{in: "dns:archive.org", opts: []MassageOpts{WithScheme}, want: "dns:archive.org"},
		{in: "dns:archive.org", opts: []MassageOpts{WithTrailingComma}, want: "dns:archive.org"},
		{in: "dns:archive.org", opts: []MassageOpts{WithScheme, WithTrailingComma}, want: "dns:archive.org"},
		{in: "whois://whois.isoc.org.il/shaveh.co.il", opts: []MassageOpts{WithScheme}, want: "whois://(il,org,isoc,whois)/shaveh.co.il"},
		{in: "whois://whois.isoc.org.il/shaveh.co.il", opts: []MassageOpts{WithTrailingComma}, want: "il,org,isoc,whois,)/shaveh.co.il"},
		{in: "whois://whois.isoc.org.il/shaveh.co.il", opts: []MassageOpts{WithTrailingComma, WithScheme}, want: "whois://(il,org,isoc,whois,)/shaveh.co.il"},
		{in: "warcinfo:foo.warc.gz", opts: []MassageOpts{WithTrailingComma}, want: "warcinfo:foo.warc.gz"},
		{in: "warcinfo:foo.warc.gz", opts: []MassageOpts{WithScheme}, want: "warcinfo:foo.warc.gz"},
		{in: "warcinfo:foo.warc.gz", opts: []MassageOpts{WithScheme, WithTrailingComma}, want: "warcinfo:foo.warc.gz"},
	}

	for _, tc := range tests {
		got := Massage(tc.in, tc.opts...)
		if got != tc.want {
			t.Errorf("surt(%q, %+v) = %q, want %q", tc.in, tc.opts, got, tc.want)
		}
	}
}
