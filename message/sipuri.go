package message

import (
	"bytes"
	"fmt"
	"strings"
)

type URIScheme []byte

var (
	SIPScheme  = URIScheme("sip:")
	SIPSScheme = URIScheme("sips:")
)

//The URIParams are defined in Section 19.1
const (
	User        = "user"
	Password    = "password"
	Host        = "host"
	Port        = "port"
	MethodParam = "method"
	MAddr       = "maddr"
	TTL         = "ttl"
	Transport   = "transport"
	LR          = "lr"
	Other       = "other"
)

type URIHeader []rune
type URIHeaderField struct {
	Name  URIHeader
	Value []rune // The value MUST be escaped following the section 19.1.2 Character Escaping Requirements
}

func (h URIHeaderField) Write(b *bytes.Buffer) {
	b.WriteString(string(h.Name))
	b.WriteByte('=')
	b.WriteString(string(h.Value))
}

//SIPURI 19.1.1 SIP-URI components sip:user:password@host:port;uri-parameters?headers
type SIPURI struct {
	Scheme  *URIScheme
	User    string
	Host    string
	Params  []ParamField
	Headers []URIHeaderField
}

func (s *SIPURI) FormatedContainsSep() bool {
	return len(s.Params) > 0 || len(s.Headers) > 0 || strings.ContainsAny(s.User, ";,")
}

func (s SIPURI) String() string {
	var b bytes.Buffer
	s.Write(&b)
	return b.String()
}

//Write method writes the SIPURI formated to a byte Buffer
func (s *SIPURI) Write(b *bytes.Buffer) error {

	if err := s.isValid(); err != nil {
		return err
	}

	//Print the scheme
	b.Write(*s.Scheme)

	//Print the user
	if len(s.User) > 0 {

		b.WriteString(s.User)
		b.WriteByte('@')
	}

	//Print the host
	b.WriteString(s.Host)

	if len(s.Params) > 0 {
		s.printParams(b)
	}

	if len(s.Headers) > 0 {
		s.printHeaders(b)
	}

	return nil
}

//isValid check if the mandatory fields are present
func (s *SIPURI) isValid() error {
	if len(s.Host) == 0 {
		return fmt.Errorf("host is mandatory")
	}

	if len(*s.Scheme) == 0 {
		return fmt.Errorf("scheme is mandatory")
	}
	return nil
}

//printParams is used by Write to format and print the Params
func (s *SIPURI) printParams(b *bytes.Buffer) {
	//Print the params
	for _, p := range s.Params {
		b.WriteByte(';')
		p.Write(b)
	}

}

//printHeaders is used by Write to format and print the headers
func (s *SIPURI) printHeaders(b *bytes.Buffer) {
	b.WriteByte('?') // ? is the separator between the params and headers
	for i, h := range s.Headers {
		if i > 0 {
			b.WriteByte('&') // & is the header separator.
		}
		h.Write(b)
	}
}
