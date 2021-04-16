package message

import (
	"bytes"
	"strings"
)

type KVP struct {
	Key   string
	Value string
}

func (p KVP) Write(b *bytes.Buffer, sep byte) {
	b.WriteString(p.Key)
	b.WriteByte(sep)
	b.WriteString(string(p.Value))
}

func ParseKVP(buf []byte, sep byte) (*KVP, bool) {
	if i := bytes.IndexByte(buf, sep); i > -1 {
		return &KVP{
			Key:   strings.TrimSpace(string(buf[:i])),
			Value: strings.TrimSpace(string(buf[i+1:])),
		}, true
	}
	return nil, false
}
