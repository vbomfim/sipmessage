package message

import (
	"bytes"
	"fmt"
	"strings"
)

type URIScheme []byte

var (
	SIPScheme     = URIScheme("sip:")
	SIPSScheme    = URIScheme("sips:")
	SIPSchemeStr  = string(SIPScheme)
	SIPSSchemeStr = string(SIPSScheme)
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

const (
	URIHeaderSep    = byte('&')
	URIHeaderSepStr = string(URIHeaderSep)
)

//SIPURI 19.1.1 SIP-URI components sip:user:password@host:port;uri-parameters?headers
type SIPURI struct {
	Scheme  *URIScheme
	User    string
	Host    string
	Params  []KVP
	Headers []KVP
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
		b.WriteByte(ParamSep)
		p.Write(b, ParamKVPSep)
	}

}

//printHeaders is used by Write to format and print the headers
func (s *SIPURI) printHeaders(b *bytes.Buffer) {
	b.WriteByte('?') // ? is the separator between the params and headers
	for i, h := range s.Headers {
		if i > 0 {
			b.WriteByte(URIHeaderSep) // & is the header separator.
		}
		h.Write(b, ParamKVPSep)
	}
}

func ParseURI(value string) (*SIPURI, error) {
	const (
		undefined = -1
	)

	var sipuri SIPURI

	err := parseScheme(value, &sipuri)
	if err != nil {
		return nil, fmt.Errorf("failed parsing uri %w", err)
	}

	var (
		iUser   = undefined
		iHost   = undefined
		iParam  = undefined
		iHeader = undefined
	)
	for i, char := range value {
		switch {
		case char == ':' && iUser == undefined:
			iUser = i + 1
		case char == '@' && iHost == undefined:
			iHost = i + 1
		case char == ';' && iHost != undefined && iParam == undefined:
			iParam = i + 1
		case char == '?' && iParam != undefined:
			iHeader = i + 1
			break
		}
	}

	if iUser == undefined || iHost == undefined {
		return nil, fmt.Errorf("failed parsing uri %s", value)
	}

	sipuri.User = value[iUser : iHost-1]

	if iParam != undefined {
		sipuri.Host = value[iHost : iParam-1]
	} else {
		if iHeader != undefined {
			sipuri.Host = value[iHost : iHeader-1]
		} else {
			sipuri.Host = value[iHost:]
		}
	}

	if iParam != undefined {
		if iHeader != undefined {
			sipuri.Params = ParseKVPs(value[iParam:iHeader-1], ParamSepStr, ParamKVPSep)
		} else {
			sipuri.Params = ParseKVPs(value[iParam:], ParamSepStr, ParamKVPSep)
		}
	}

	if iHeader != undefined {
		sipuri.Headers = ParseKVPs(value[iHeader:], URIHeaderSepStr, ParamKVPSep)
	}

	return &sipuri, nil
}

func parseScheme(value string, sipuri *SIPURI) error {
	switch {
	case value[:4] == SIPSchemeStr:
		sipuri.Scheme = &SIPScheme
	case value[:5] == SIPSSchemeStr:
		sipuri.Scheme = &SIPSScheme
	default:
		return fmt.Errorf("invalid scheme")
	}
	return nil
}
