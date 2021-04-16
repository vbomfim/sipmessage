package message

import (
	"bytes"
	"strconv"
)

type CSeq struct {
	Seq    uint32
	Method Method
}

func (_ CSeq) Tag() string {
	return "CSeq"
}

func (c CSeq) String() string {
	var b bytes.Buffer
	b.WriteString(strconv.FormatUint(uint64(c.Seq), 10))
	b.WriteByte(SP)
	b.WriteString(c.Method)

	return b.String()
}
