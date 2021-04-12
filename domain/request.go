package domain

// https://tools.ietf.org/html/rfc3261
// 7.1 - Method: This specification defines six methods:
// REGISTER for registering contact information,
// INVITE, ACK, and CANCEL for setting up sessions,
// BYE for terminating sessions, and
// OPTIONS for querying servers about their capabilities.
// SIP extensions, documented in standards track RFCs, may define
// additional methods.
var Methods = []string{
	"REGISTER",
	"INVITE",
	"ACK",
	"CANCEL",
	"BYE",
	"OPTIONS",

	"SUBSCRIBE",
	"NOTIFY",
	"REFER",
	"INFO",
}

type Request struct {
	RawContent  string
	Method      string
	RequestURI  string
	SIPVersion  string
	To          string
	From        string
	CSeq        string
	CallID      string
	MaxForwards string
	Via         []string
}
