package message

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

//Response type holds the data from a User Agent Server(UAS) in response of a Request.
type Response struct {
	StatusCode   string
	ReasonPhrase string
	Headers      []Header
	Body         []byte
}

func (r Response) String() string {
	var b bytes.Buffer
	if r.Write(&b) != nil {
		return "INVALIDRESPONSE"
	}
	return b.String()
}

func (r Response) Write(w io.Writer) error {
	var b bytes.Buffer
	b.WriteString(Version)
	b.WriteByte(SP)
	b.WriteString(r.StatusCode)
	b.WriteByte(SP)
	b.WriteString(r.ReasonPhrase)
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

//ParseResponse parses a string to a Response instance
func ParseResponse(rawData string) (*Response, error) {
	var (
		resp Response
		err  error
	)

	i := strings.Index(rawData, strCRLF) // index of the first EOL
	err = parseResponseLine(rawData[:i], &resp)
	if err != nil {
		return nil, fmt.Errorf("invalid argument %w", err)
	}

	emptyLineIndex := strings.Index(rawData, strDoubleCRLF) // index of an empty line after the headers.

	headers, err := ParseHeaders(rawData[i+2 : emptyLineIndex]) // +2 means CRLF size
	if err != nil {
		return nil, fmt.Errorf("invalid argument %w", err)
	}
	resp.Headers = headers

	//TODO: Body parse?

	return &resp, nil
}

//parseRequestLine parses the Request-Line
func parseResponseLine(line string, resp *Response) error {
	const (
		fieldsN           = 3
		fieldVersion      = 0
		fieldStatusCode   = 1
		fieldReasonPhrase = 2
	)
	fields := strings.Fields(line)

	// Confirms the number of fields
	if len(fields) != fieldsN {
		return fmt.Errorf("mal-formatted response-line %s", line)
	}

	// Confirms the version is compatible.
	if Version != fields[fieldVersion] {
		return fmt.Errorf("incompatible sip version %s", line)
	}

	resp.StatusCode = fields[fieldStatusCode]
	resp.ReasonPhrase = fields[fieldReasonPhrase]

	return nil
}
