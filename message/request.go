package message

import (
	"bytes"
	"io"
)

//Version hold the SIP version string used on Request Lines
const Version = "SIP/2.0"

//Method https://tools.ietf.org/html/rfc3261
//The method is the primary function that a REQUEST is meant
//to invoke on a server.  The method is carried in the request
//message itself.  Example methods are INVITE and BYE.
type Method = []byte

// https://tools.ietf.org/html/rfc3261
// 7.1 - Method: This specification defines six methods:
// REGISTER for registering contact information,
// INVITE, ACK, and CANCEL for setting up sessions,
// BYE for terminating sessions, and
// OPTIONS for querying servers about their capabilities.
// SIP extensions, documented in standards track RFCs, may define
// additional methods.
var (
	REGISTER  = Method("REGISTER")
	INVITE    = Method("INVITE")
	ACK       = Method("ACK")
	CANCEL    = Method("CANCEL")
	BYE       = Method("BYE")
	OPTIONS   = Method("OPTIONS")
	SUBSCRIBE = Method("SUBSCRIBE")
	NOTIFY    = Method("NOTIFY")
	REFER     = Method("REFER")
	INFO      = Method("INFO")
)

var (
	//SP means space. This is the separator between the Request Line fields.
	SP = byte(' ')
)

//Request type holds the data that identifies a request from a User Agent Client(UAC).
//The mandatory headers are members of this type. The remaining headers and the body
//can be extracted from the RawContent attribute through helper methods.
//This approach avoids unnecessary parsing.
type Request struct {
	Method     Method
	RequestURI string //RFC 3261 Section 8.1.1.1
	Headers    []Header
	Body       []byte
}

func (r Request) Write(w io.Writer) error {
	var b bytes.Buffer
	// RequestLine
	b.Write([]byte(r.Method))
	b.WriteByte(SP)
	b.WriteString(string(r.RequestURI))
	b.WriteByte(SP)
	b.Write([]byte(Version))
	b.Write(CRLF)

	//Headers
	for _, h := range r.Headers {
		WriteHeader(&b, h)
	}

	//Header and Body separation of 1 line
	b.Write(CRLF)

	if r.Body != nil {
		b.Write([]byte(string(r.Body)))
	}
	_, err := w.Write(b.Bytes())
	return err
}
