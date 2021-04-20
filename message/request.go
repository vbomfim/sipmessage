package message

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

//Version hold the SIP version string used on Request Lines
const Version = "SIP/2.0"

//Method https://tools.ietf.org/html/rfc3261
//The method is the primary function that a REQUEST is meant
//to invoke on a server.  The method is carried in the request
//message itself.  Example methods are INVITE and BYE.
type Method = string

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

	Methods = []Method{REGISTER, INVITE, ACK, CANCEL, BYE, OPTIONS, SUBSCRIBE, NOTIFY, REFER, INFO}
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

func (r Request) String() string {
	var b bytes.Buffer
	r.Write(&b)
	return b.String()
}

func (r Request) Write(w io.Writer) error {
	var b bytes.Buffer
	// RequestLine
	b.Write([]byte(r.Method))
	b.WriteByte(SP)
	b.WriteString(string(r.RequestURI))
	b.WriteByte(SP)
	b.WriteString(Version)
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

//ParseRequest parses a string to a Request instance
func ParseRequest(rawData string) (*Request, error) {
	var (
		req Request
		err error
	)

	i := strings.Index(rawData, strCRLF) // index of the first EOL
	err = parseRequestLine(rawData[:i], &req)
	if err != nil {
		return nil, fmt.Errorf("invalid argument %w", err)
	}

	emptyLineIndex := strings.Index(rawData, strDoubleCRLF) // index of an empty line after the headers.

	headers, err := ParseHeaders(rawData[i+2 : emptyLineIndex]) // +2 means CRLF size
	if err != nil {
		return nil, fmt.Errorf("invalid argument %w", err)
	}
	req.Headers = headers

	//TODO: Body parse?

	return &req, nil
}

//parseRequestLine parses the Request-Line
func parseRequestLine(line string, req *Request) error {
	const (
		fieldsN         = 3
		fieldMethod     = 0
		fieldRequestURI = 1
		fieldVersion    = 2
	)
	fields := strings.Fields(line)

	// Confirms the number of fields
	if len(fields) != fieldsN {
		return fmt.Errorf("mal-formatted request-line %s", line)
	}

	// Parsing the Request-Line Method field
	if method, OK := parseMethod(fields[fieldMethod]); OK {
		req.Method = *method
	} else {
		return fmt.Errorf("mal-formatted Method in the request-line %s", line)
	}
	// Assign the requestURI as string. It doesn't parse.
	req.RequestURI = fields[fieldRequestURI]

	// Confirms the version is compatible.
	if Version != fields[fieldVersion] {
		return fmt.Errorf("incompatible sip version %s", line)
	}

	return nil
}

//parseMethod verify if the value is a valid Method and returns the corresponding Method instance
func parseMethod(value string) (*Method, bool) {
	for _, m := range Methods {
		if m == Method(value) {
			return &m, true
		}
	}
	return nil, false
}
