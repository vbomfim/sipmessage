package message

import "bytes"

type ParamField struct {
	Name  string
	Value string
}

func (p ParamField) Write(b *bytes.Buffer) {
	b.WriteString(p.Name)
	b.WriteByte('=')
	b.WriteString(string(p.Value))
}
