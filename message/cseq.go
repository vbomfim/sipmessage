package message

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

//CSeq represents the CSeq SIP Header
type CSeq struct {
	Seq    uint32
	Method Method
}

//Header returns the header tag
func (CSeq) Header() string {
	return "CSeq"
}

//CHeader returns the compact header tag
func (CSeq) CHeader() string {
	return "CSeq"
}

func (c CSeq) String() string {
	var b bytes.Buffer
	b.WriteString(strconv.FormatUint(uint64(c.Seq), 10))
	b.WriteByte(SP)
	b.WriteString(c.Method)

	return b.String()
}

//ParseCSeq parses a string value to a CSeq instance
func ParseCSeq(value string) (*CSeq, error) {
	fields := strings.Fields(value)
	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid CSeq format - %s", value)

	}
	const (
		fieldSeq    = 0
		fieldMethod = 1
	)

	// Parsing the Sequence field
	v, err := strconv.ParseUint(fields[fieldSeq], 10, 32)
	if err != nil {
		return nil, err
	}
	cseq := CSeq{
		Seq: uint32(v),
	}

	// Parsing the Method field
	if method, OK := parseMethod(fields[fieldMethod]); OK {
		cseq.Method = *method
	} else {
		return nil, fmt.Errorf("mal-formatted Method parsing CSeq - %s", fields[fieldMethod])
	}

	return &cseq, nil
}
