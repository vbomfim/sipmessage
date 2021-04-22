package message

import (
	"bytes"
	"strings"
)

//KVP is a key value pair type used to define ordered Params in KPV slices.
type KVP struct {
	Key   string
	Value string
}

func (p KVP) Write(b *bytes.Buffer, sep byte) {
	b.WriteString(p.Key)
	b.WriteByte(sep)
	b.WriteString(string(p.Value))
}

//WriteKVPs writes to a byte buffer a list of KVPs
func WriteKVPs(b *bytes.Buffer, kvps []KVP, KVPsSep byte, sep byte) {
	//Print the params
	for _, p := range kvps {
		b.WriteByte(KVPsSep)
		p.Write(b, sep)
	}
}

//ParseKVP parses a byte slice to a KVP instance
//The second parameter sep is the separator between the key and the value
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
func ParseKVPs(value, KVPsSep string, sep byte) []KVP {
	var result = []KVP{}
	parts := strings.Split(value, KVPsSep)
	for _, part := range parts {
		if kvp, OK := ParseKVP([]byte(part), sep); OK {
			result = append(result, *kvp)
		}
	}
	return result
}
