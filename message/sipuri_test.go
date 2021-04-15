package message_test

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/vbomfim/sipmessage/message"
)

func TestURIFormat(t *testing.T) {
	const want = "sip:alice@atlanta.com;transport=TCP?subject=project%20x&priority=urgent&to=alice%40atlanta.com"
	user := "alice"
	host := "atlanta.com"
	params := []message.ParamField{
		{Name: message.Transport, Value: "TCP"},
	}
	uriHeaders := []message.URIHeaderField{
		{Name: []rune("subject"), Value: []rune(url.PathEscape("project x"))},
		{Name: []rune("priority"), Value: []rune("urgent")},
		{Name: []rune("to"), Value: []rune(url.QueryEscape("alice@atlanta.com"))},
	}

	uri := message.SIPURI{
		Scheme:  &message.SIPScheme,
		User:    user,
		Host:    host,
		Params:  params,
		Headers: uriHeaders,
	}

	var b bytes.Buffer
	err := uri.Write(&b)
	if err != nil {
		t.Fatalf("failed on write method - %v", err)
	}

	got := b.String()
	if got != want {
		t.Fatalf("want different from got\nwant:\n%s\ngot:\n%s", want, got)
	}

}
