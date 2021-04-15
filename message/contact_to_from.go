package message

import (
	"bytes"
	"strings"
)

//To
type To Contact

func (_ To) Tag() string {
	return "To"
}

func (c To) String() string {
	return printInContactFormat(c.DisplayName, c.SIPURI, c.Params)
}

//From
type From Contact

func (_ From) Tag() string {
	return "From"
}

func (c From) String() string {
	return printInContactFormat(c.DisplayName, c.SIPURI, c.Params)
}

//Contact
type Contact struct {
	DisplayName string
	SIPURI      SIPURI
	Params      []ParamField
}

func (_ *Contact) Tag() string {
	return "Contact"
}

func (c Contact) String() string {
	return printInContactFormat(c.DisplayName, c.SIPURI, c.Params)
}

func printInContactFormat(displayName string, SIPURI SIPURI, params []ParamField) string {
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
			b.WriteByte(';')
			p.Write(&b)
		}
	}

	return b.String()
}
