package message

import "bytes"

//Header is a type used by the SIP Headers. []byte is a convenient way to parse and write.
type Header []byte

//Exports the Headers defined by https://tools.ietf.org/html/rfc3261
var (
	//To - section 8.1.1.2
	To = Header("TO")
	//From - section 8.1.1.3
	From = Header("FROM")
	//CSeq - section 8.1.1.5
	CSeq = Header("CSEQ")
	//CallID - section 8.1.1.4
	CallID = Header("CALL-ID")
	//MaxForwards - section 8.1.1.6
	MaxForwards = Header("MAX-FORWARDS")
	//Via - section 8.1.1.7
	Via = Header("VIA")
	//Contact - section 8.1.1.8
	Contact            = Header("CONTACT")
	ContentLength      = Header("CONTENT-LENGTH")
	ContentDisposition = Header("CONTENT-DISPOSITION")
	ContentEncoding    = Header("CONTENT-ENCODING")
	ContentLanguage    = Header("CONTENT-LANGUAGE")
	Accept             = Header("ACCEPT")
	AcceptEncoding     = Header("ACCEPT-ENCODING")
	AcceptLanguage     = Header("ACCEPT-LANGUAGE")
	AlertInfo          = Header("ALERT-INFO")
	Allow              = Header("ALLOW")
	AuthenticationInfo = Header("AUTHENTICATION-INFO")
	Authorization      = Header("AUTHORIZATION")
	CallInfo           = Header("CALL-INFO")
)

var (
	//HeaderSep is the separator used between the Header and its value when printed
	HeaderSep = []byte(": ")
	//CRLF defines the end of line.
	CRLF = []byte{'\r', '\n'}
)

//HeaderField is a helper type to print a Header and its value properly - section 7.3.1
type HeaderField struct {
	Header Header
	Value  []rune
}

//Write method writes the HeaderField following the section 7.3.1
func (h HeaderField) Write(b *bytes.Buffer) {
	b.Write(h.Header)
	b.Write(HeaderSep)
	b.Write([]byte(string(h.Value)))
	b.Write(CRLF)
}
