package message

import "bytes"

//Header is a type used by the SIP Headers. []byte is a convenient way to parse and write.
type Header interface {
	Tag() string
	String() string
}

var (
	//HeaderSep is the separator used between the Header and its value when printed
	HeaderSep = []byte(": ")
	//CRLF defines the end of line.
	CRLF = []byte{'\r', '\n'}
)

//Write method writes the HeaderField following the section 7.3.1
func WriteHeader(b *bytes.Buffer, h Header) {
	b.WriteString(h.Tag())
	b.Write(HeaderSep)
	b.WriteString(h.String())
	b.Write(CRLF)
}
