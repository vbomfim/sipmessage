package message

//Response type holds the data from a User Agent Server(UAS) in response of a Request.
type Response struct {
	RawContent   string
	SIPVersion   string
	StatusCode   string
	ReasonPhrase string
	Headers      []string
}
