package surt

type MassageOpts int

const (
	// WithScheme indicates that the massaged URL should include the scheme. e.g.: "http://example.com" -> "http://(com,example)/"
	WithScheme MassageOpts = iota
	// WithTrailingComma indicates that the massaged URL should end with a comma. e.g.: "http://example.com" -> "http://(com,example,)/"
	WithTrailingComma
	//WithIACanonicalization indicates that the massaged URL should be canonicalized according to IA rules.
	WithIACanonicalization
	//NoHostMassage indicates that the host shouldn't be massaged (e.g., stripping "www.").
	NoHostMassage
)
