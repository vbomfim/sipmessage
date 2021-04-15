package message_test

import (
	"testing"

	"github.com/vbomfim/sipmessage/message"
)

func TestContactHeader(t *testing.T) {
	t.Run("Testing Contact Header String", func(t *testing.T) {
		want := "Bob <sip:bob@biloxy.example.com>"
		contact := message.Contact{DisplayName: "Bob", SIPURI: message.SIPURI{Scheme: &message.SIPScheme, User: "bob", Host: "biloxy.example.com"}}
		got := contact.String()
		if want != got {
			t.Fatalf("contact header mal-formatted\nwant:\n%s\ngot:\n%s", want, got)
		}
	})
}

func TestToHeader(t *testing.T) {
	t.Run("Testing To Header String", func(t *testing.T) {
		want := "Bob <sip:bob@biloxy.example.com>"
		to := message.To{DisplayName: "Bob", SIPURI: message.SIPURI{Scheme: &message.SIPScheme, User: "bob", Host: "biloxy.example.com"}}
		got := to.String()
		if want != got {
			t.Fatalf("To header mal-formatted\nwant:\n%s\ngot:\n%s", want, got)
		}
	})
	t.Run("Testing the TAG", func(t *testing.T) {
		want := "To"
		to := message.To{}
		got := to.Tag()
		if want != got {
			t.Fatalf("wrong TAG\nwant:\n%s\ngot:\n%s", want, got)
		}
	})
}
