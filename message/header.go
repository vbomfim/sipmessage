package message

import (
	"bufio"
	"bytes"
	"strings"
)

const (
	ParamSep    = byte(';')
	ParamSepStr = string(ParamSep)
	ParamKVPSep = byte('=')
)

var (
	//HeaderSep is the separator used between the Header and its value when printed
	HeaderSep = byte(':')
	//CRLF defines the end of line.
	CRLF          = []byte{'\r', '\n'}
	strCRLF       = string(CRLF)
	strDoubleCRLF = string([]byte{'\r', '\n', '\r', '\n'})
)

//Header is a type used by the SIP Headers. []byte is a convenient way to parse and write.
type Header interface {
	Tag() string
	CTag() string // Tag in the compact form
	String() string
}

//Write method writes the HeaderField following the section 7.3.1
func WriteHeader(b *bytes.Buffer, h Header) {
	b.WriteString(h.Tag())
	b.WriteByte(HeaderSep)
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

func ParseHeader(b []byte) (Header, error) {

	if kvp, OK := ParseKVP(b, HeaderSep); OK {
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
		}
	}

	return nil, nil
}

func isSpaceOrTab(char byte) bool {
	return char == ' ' || char == '\t'
}
