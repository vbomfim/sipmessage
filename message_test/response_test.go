package message_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vbomfim/sipmessage/message"
)

func TestResponse(t *testing.T) {
	respF := getRegisterMessageResponseSample()
	resp, _ := message.ParseResponse(respF)
	fmt.Println(resp)
}

func getRegisterMessageResponseSample() string {
	var sb strings.Builder
	sb.WriteString("SIP/2.0 200 OK\r\n")
	sb.WriteString("Via: SIP / 2.0 / TLS   client.biloxi.example.com:  5061   ;branch=z9hG4bKnashd92;received=192.0.2.201\r\n")
	sb.WriteString("From: Bob <sips:bob@biloxi.example.com>;tag=ja743ks76zlflH\r\n")
	sb.WriteString("To: Bob <sips:bob@biloxi.example.com>;tag=37GkEhwl6\r\n")
	sb.WriteString("Call-ID: 1j9FpLxk3uxtm8tn@biloxi.example.com\r\n")
	sb.WriteString("CSeq: 2 REGISTER\r\n")
	sb.WriteString("Contact: <sips:bob@client.biloxi.example.com>;expires=3600\r\n")
	sb.WriteString("Content-Length: 0\r\n")
	sb.WriteString("\r\n")
	return sb.String()
}
