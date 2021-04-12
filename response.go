package domain

type Response struct {
	RawContent   string
	SIPVersion   string
	StatusCode   string
	ReasonPhrase string
	Headers      []string
}
