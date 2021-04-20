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

//ParseKVPs parses a string with KVPs separated by a token and returns a slice of KVPs
//	Examples:
//	value: "param1=1234;param2=5678" KVPsSep: ";"  KVPSep: '='
//	value: "param1=abcd&param2=efgh" KVPsSep: "&"  KVPSep: '='
func ParseKVPs(value, KVPsSep string, KVPSep byte) []KVP {
	var result = []KVP{}
	parts := strings.Split(value, KVPsSep)
	for _, part := range parts {
		if kvp, OK := ParseKVP([]byte(part), KVPSep); OK {
			result = append(result, *kvp)
		}
	}
	return result
}
