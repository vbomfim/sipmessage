package message_test

import (
	"fmt"
	"testing"

	"github.com/vbomfim/sipmessage/message"
	"github.com/vbomfim/sipmessage/message/param"
)

func TestParseVia(t *testing.T) {
	data := fmt.Sprintf("Via: SIP / 2.0 / UDP   \t first.example.com: \t\t\t4000;ttl=16 ;maddr=224.2.0.1 ;branch=z9hG4bKa7c6a8dlze.1")
	via, _ := message.ParseVia(data)
	fmt.Println(data)
	fmt.Println(via)
}

func TestPrintVia(t *testing.T) {
	want := fmt.Sprintf("SIP/2.0/UDP first.example.com:4000;ttl=16;maddr=224.2.0.1;branch=z9hG4bKa7c6a8dlze.1\r\n")
	via := message.Via{
		Transport: message.UDP,
		Host:      "first.example.com:4000",
		Params: []message.KVP{
			{Key: param.TTL, Value: "16"},
			{Key: param.MAddr, Value: "224.2.0.1"},
			{Key: param.Branch, Value: "z9hG4bKa7c6a8dlze.1"},
		},
	}
	got := via.String()
	fmt.Println(got)
	if want != got {
		t.Fatalf("want: \n%s\n but got:\n%s\n", want, got)
	}

}
