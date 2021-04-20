package message_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vbomfim/sipmessage/message"
)

func TestParseHeaders(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("Max-Forwards  : 105 \r\n")
	sb.WriteString("From: Bob <sips:bob@biloxi.example.com>;tag=a73kszlfl\r\n")
	sb.WriteString("To: \"A. G. Bell\" <sips:agb@bell-telephone.com>\r\n")
	sb.WriteString("Call-ID: 1j9FpLxk3uxtm8tn@biloxi.example.com,\r\n")
	sb.WriteString(" \"ai caramba\"\r\n")
	sb.WriteString("CSeq  : 1 REGISTER  \r\n")
	sb.WriteString("Contact: Bob <sips:bob@biloxi.example.com>\r\n")
	sb.WriteString("Content-Length: 0\r\n")
	sb.WriteString("\r\n")
	headers, _ := message.ParseHeaders(sb.String())
	for _, header := range headers {
		fmt.Println(header.Tag(), "->", header.String())
	}

}
