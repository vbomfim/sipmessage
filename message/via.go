package message

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

//Defining errors
var (
	ErrInvalidTransport       = errors.New("message.Via: Invalid Transport")
	ErrInvalidViaHeaderFormat = errors.New("message.Via: Invalid Via Header Format")
)

//Transport represents the transport protocol
type Transport = string

//Transport constants definitions
const (
	UDP  Transport = "UDP"
	TCP  Transport = "TCP"
	TLS  Transport = "TLS"
	SCTP Transport = "SCTP"
)

//Transports holds a list of valid Transports
var Transports = []Transport{UDP, TCP, TLS, SCTP}

//Via represents the Via SIP header
type Via struct {
	Transport Transport
	Host      string
	Params    []KVP
}

func (Via) Header() string {
	return "Via"
}

func (Via) CHeader() string {
	return "v"
}

func (v Via) String() string {
	var b bytes.Buffer
	// RequestLine
	b.WriteString(Version)
	b.WriteByte('/')
	b.WriteString(string(v.Transport))
	b.WriteByte(SP)
	b.WriteString(v.Host)
	if len(v.Params) > 0 {
		WriteKVPs(&b, v.Params, paramSep, paramKVPSep)
	}
	return b.String()
}

//ParseVia parses a string value to the Via type
func ParseVia(value string) (*Via, error) {
	const (
		undefined = -1
	)

	var via Via

	var (
		iTransport = undefined
		iParam     = undefined
	)
	// Look for the separators
Loop:
	for i, char := range value {
		switch char {
		case '/':
			iTransport = i + 1
		case ';':
			iParam = i + 1
			break Loop
		}
	}
	if iTransport == undefined {
		return nil, ErrInvalidViaHeaderFormat
	}

	var transpHost string
	if iParam == undefined {
		transpHost = strings.TrimSpace(value[iTransport:])
	} else {
		transpHost = strings.TrimSpace(value[iTransport : iParam-1])
		via.Params = ParseKVPs(value[iParam:], paramSepStr, paramKVPSep)
	}

	iHost := strings.Index(transpHost, string(SP))
	if iHost >= 0 {
		err := parseTransport(transpHost[:iHost], &via)
		if err != nil {
			return nil, err
		}
		// removes spaces and tabs between the host and port
		spacePattern := regexp.MustCompile(`\s+|\t+`)
		via.Host = spacePattern.ReplaceAllString(transpHost[iHost:], "")
	}

	return &via, nil
}

func parseTransport(value string, via *Via) error {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case UDP:
		via.Transport = UDP
	case TCP:
		via.Transport = TCP
	case TLS:
		via.Transport = TLS
	case SCTP:
		via.Transport = SCTP
	default:
		return ErrInvalidTransport
	}
	return nil
}
