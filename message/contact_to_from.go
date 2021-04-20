package message

import (
	"bytes"
	"fmt"
	"strings"
)

//To
type To struct {
	DisplayName string
	SIPURI      SIPURI
	Params      []KVP
}

func (_ To) Tag() string {
	return "To"
}

func (_ To) CTag() string {
	return "t"
}

func (c To) String() string {
	return printInContactFormat(c.DisplayName, c.SIPURI, c.Params)
}

//From
type From struct {
	DisplayName string
	SIPURI      SIPURI
	Params      []KVP
}

func (_ From) Tag() string {
	return "From"
}

func (_ From) CTag() string {
	return "f"
}

func (c From) String() string {
	return printInContactFormat(c.DisplayName, c.SIPURI, c.Params)
}

//Contact
type Contact struct {
	DisplayName string
	SIPURI      SIPURI
	Params      []KVP
}

func (_ Contact) Tag() string {
	return "Contact"
}

func (_ Contact) CTag() string {
	return "m"
}

func (c Contact) String() string {
	return printInContactFormat(c.DisplayName, c.SIPURI, c.Params)
}

func printInContactFormat(displayName string, SIPURI SIPURI, params []KVP) string {
	var b bytes.Buffer

	needAngleBrackets := displayName != "" || len(params) > 0 || SIPURI.FormatedContainsSep()

	if displayName != "" {
		containSpace := strings.Contains(displayName, " ")
		if containSpace {
			b.WriteByte('"')
		}
		b.WriteString(displayName)
		if containSpace {
			b.WriteByte('"')
		}
		b.WriteByte(' ')
	}

	if needAngleBrackets {
		b.WriteByte('<')
	}

	SIPURI.Write(&b)

	if needAngleBrackets {
		b.WriteByte('>')
	}

	if len(params) > 0 {
		//Print the params
		for _, p := range params {
			b.WriteByte(ParamSep)
			p.Write(&b, ParamKVPSep)
		}
	}

	return b.String()
}

type NewContactHeaderFunc func(displayName string, SIPURI *SIPURI, params []KVP) (Header, error)

func NewTo(displayName string, sipuri *SIPURI, params []KVP) (Header, error) {
	if sipuri == nil {
		return nil, fmt.Errorf("sipuri is mandatory")
	}
	return To{
		DisplayName: displayName,
		SIPURI:      *sipuri,
		Params:      params,
	}, nil
}
func NewFrom(displayName string, sipuri *SIPURI, params []KVP) (Header, error) {
	if sipuri == nil {
		return nil, fmt.Errorf("sipuri is mandatory")
	}
	return From{
		DisplayName: displayName,
		SIPURI:      *sipuri,
		Params:      params,
	}, nil
}

func NewContact(displayName string, sipuri *SIPURI, params []KVP) (Header, error) {
	if sipuri == nil {
		return nil, fmt.Errorf("sipuri is mandatory")
	}
	return Contact{
		DisplayName: displayName,
		SIPURI:      *sipuri,
		Params:      params,
	}, nil
}

/*ParseContact
Contact        =  ("Contact" / "m" ) HCOLON
                  ( STAR / (contact-param *(COMMA contact-param)))
contact-param  =  (name-addr / addr-spec) *(SEMI contact-params)
name-addr      =  [ display-name ] LAQUOT addr-spec RAQUOT
addr-spec      =  SIP-URI / SIPS-URI / absoluteURI
display-name   =  *(token LWS)/ quoted-string

contact-params     =  c-p-q / c-p-expires
                      / contact-extension
c-p-q              =  "q" EQUAL qvalue
c-p-expires        =  "expires" EQUAL delta-seconds
contact-extension  =  generic-param
*/
func ParseContact(value string, new NewContactHeaderFunc) (Header, error) {

	/*
		Option 1: *
		Option 2: SIP-URI / SIPS-URI / absoluteURI;Contact Params
		Option 3: <SIP-URI with Params>;Contact Params
		Option 3: TOKEN <OPTION 2>
		Option 4: "any"<OPTION 2>
	*/
	const (
		undefined = -1
	)
	var (
		iOpenQuotes         = undefined
		iCloseQuotes        = undefined
		iOpenAngleBrackets  = undefined
		iCloseAngleBrackets = undefined
		iParamSep           = undefined
	)
	for i, char := range value {
		switch {
		case char == '"' && iOpenQuotes == undefined:
			iOpenQuotes = i
		case char == '"' && iOpenQuotes != undefined:
			iCloseQuotes = i
		case char == '<':
			iOpenAngleBrackets = i
		case char == '>':
			iCloseAngleBrackets = i
		case char == ';' && (iOpenAngleBrackets == undefined || iOpenAngleBrackets != undefined && iCloseAngleBrackets != undefined && iParamSep == undefined):
			iParamSep = i
		}
	}

	if (iOpenQuotes == undefined) != (iCloseQuotes == undefined) || (iOpenAngleBrackets == undefined) != (iCloseAngleBrackets == undefined) {
		return nil, fmt.Errorf("invalid format")
	}
	var (
		displayName string
		params      []KVP
		sipuri      *SIPURI
		err         error
	)

	if hasDisplayName := iOpenAngleBrackets > 0; hasDisplayName {
		if iOpenQuotes == undefined {
			displayName = strings.TrimSpace(value[:iOpenAngleBrackets])
		} else {
			displayName = strings.TrimSpace(value[iOpenQuotes+1 : iCloseQuotes])
		}
	}

	hasParams := iParamSep != undefined

	if iOpenAngleBrackets != undefined {
		sipuri, err = ParseURI(value[iOpenAngleBrackets+1 : iCloseAngleBrackets])
	} else {
		if hasParams {
			sipuri, err = ParseURI(value[:iParamSep])
		} else {
			sipuri, err = ParseURI(value)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("invalid URI format")
	}

	if hasParams {
		params = ParseKVPs(value[iParamSep+1:], ParamSepStr, ParamKVPSep)
	}

	return new(displayName, sipuri, params)
}
