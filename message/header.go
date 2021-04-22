package message

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
)

const (
	paramSep    = byte(';')
	paramSepStr = string(paramSep)
	paramKVPSep = byte('=')
)

//Defining errors
var (
	ErrUnknownHeader = errors.New("message.Header: Unknown header")
)
var (
	//headerSep is the separator used between the Header and its value when printed
	headerSep = byte(':')
	//CRLF defines the end of line.
	CRLF          = []byte{'\r', '\n'}
	strCRLF       = string(CRLF)
	strDoubleCRLF = string([]byte{'\r', '\n', '\r', '\n'})
)

//Header is a type used by the SIP Headers. []byte is a convenient way to parse and write.
type Header interface {
	Header() string
	CHeader() string // Tag in the compact form
	String() string
}

//WriteHeader method writes the HeaderField following the section 7.3.1
func WriteHeader(b *bytes.Buffer, h Header) {
	b.WriteString(h.Header())
	b.WriteByte(headerSep)
	b.WriteByte(SP)
	b.WriteString(h.String())
	b.Write(CRLF)
}

//ParseHeaders follows the section 7.3.1 Header Field Format
//It receives a string with one or more headers
func ParseHeaders(line string) ([]Header, error) {
	var headers []Header
	var b bytes.Buffer

	r := strings.NewReader(line)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		if len(line) == 0 {
			continue
		}
		if isSpaceOrTab(line[0]) { // A space or a tab at the begining means that the content belongs to the previous line
			b.WriteByte(SP)
			b.WriteString(strings.TrimSpace(line))
			continue
		}
		if b.Len() > 0 {
			header, err := ParseHeader(b.Bytes())
			if err == nil {
				headers = append(headers, header)
			}
		}
		b.Reset()
		b.WriteString(line)
	}

	if b.Len() > 0 {
		header, err := ParseHeader(b.Bytes())
		if err == nil {
			headers = append(headers, header)
		}
	}
	return headers, nil
}

//ParseHeader parses a byte slice value to a Header instance
func ParseHeader(b []byte) (Header, error) {

	if kvp, OK := ParseKVP(b, headerSep); OK {
		switch strings.ToLower(kvp.Key) {
		case "max-forwards":
			return ParseMaxForwards(kvp.Value)
		case "from", "f":
			return ParseContact(kvp.Value, NewFrom)
		case "to", "t":
			return ParseContact(kvp.Value, NewTo)
		case "contact", "m":
			return ParseContact(kvp.Value, NewContact)
		case "call-id", "i":
			return ParseCallID(kvp.Value)
		case "content-length", "l":
			return ParseContentLength(kvp.Value)
		case "cseq":
			return ParseCSeq(kvp.Value)
		case "via", "v":
			return ParseVia(kvp.Value)
		}
	}

	return nil, ErrUnknownHeader
}

func isSpaceOrTab(char byte) bool {
	return char == ' ' || char == '\t'
}
