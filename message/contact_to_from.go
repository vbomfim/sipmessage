package message

import (
	"bytes"
	"fmt"
	"strings"
)

//To represents the To SIP Header
type To struct {
	DisplayName string
	SIPURI      SIPURI
	Params      []KVP
}

//Header returns the header tag
func (To) Header() string {
	return "To"
}

//CHeader returns the compact header tag
func (To) CHeader() string {
	return "t"
}

func (c To) String() string {
	return printContactFormat(c.DisplayName, c.SIPURI, c.Params)
}

//From represents the From SIP Header
type From struct {
	DisplayName string
	SIPURI      SIPURI
	Params      []KVP
}

//Header returns the header tag
func (From) Header() string {
	return "From"
}

//CHeader returns the compact header tag
func (From) CHeader() string {
	return "f"
}

func (c From) String() string {
	return printContactFormat(c.DisplayName, c.SIPURI, c.Params)
}

//Contact represents the Contact SIP Header
type Contact struct {
	DisplayName string
	SIPURI      SIPURI
	Params      []KVP
}

//Header returns the header tag
func (Contact) Header() string {
	return "Contact"
}

//CHeader returns the compact header tag
func (Contact) CHeader() string {
	return "m"
}

func (c Contact) String() string {
	return printContactFormat(c.DisplayName, c.SIPURI, c.Params)
}

func printContactFormat(displayName string, SIPURI SIPURI, params []KVP) string {
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

	err := SIPURI.Write(&b)
	if err != nil {
		b.WriteString("INVALIDURI")
	}

	if needAngleBrackets {
		b.WriteByte('>')
	}

	if len(params) > 0 {
		//Print the params
		for _, p := range params {
			b.WriteByte(paramSep)
			p.Write(&b, paramKVPSep)
		}
	}

	return b.String()
}

//NewContactHeaderFunc is a function type to be implemented by Constructors of Contact style types like From, To and Contact
type NewContactHeaderFunc func(displayName string, SIPURI *SIPURI, params []KVP) (Header, error)

//NewTo is a constructor of the To Header
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

//NewFrom is a constructor of the From Header
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

//NewContact is a constructor of the Contact Header
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

//ParseContact parses a string value to a Contact instance and returns as a Header interface
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
		params = ParseKVPs(value[iParamSep+1:], paramSepStr, paramKVPSep)
	}

	return new(displayName, sipuri, params)
}
