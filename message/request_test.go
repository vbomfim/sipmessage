package message_test

import (
	"strings"
	"testing"

	"github.com/vbomfim/sipmessage/message"
)

type OutputStub struct {
	value []byte
}

func (o *OutputStub) Write(p []byte) (n int, err error) {
	o.value = p
	return len(p), nil
}
func TestMessageFormat(t *testing.T) {
	t.Run("Create new Register message", func(t *testing.T) {

		want := getRegisterMessageSample()
		reqURI := message.SIPURI{
			Scheme: &message.SIPSScheme,
			Host:   "ss2.biloxi.example.com",
		}

		contact := message.Contact{DisplayName: "Bob", SIPURI: message.SIPURI{Scheme: &message.SIPSScheme, User: "bob", Host: "biloxi.example.com"}}
		from := message.From{DisplayName: contact.DisplayName, SIPURI: contact.SIPURI}
		from.Params = []message.ParamField{{Name: "tag", Value: "a73kszlfl"}}
		to := message.To{DisplayName: "A. G. Bell", SIPURI: message.SIPURI{Scheme: &message.SIPSScheme, User: "agb", Host: "bell-telephone.com"}}
		callid := message.CallID("1j9FpLxk3uxtm8tn@biloxi.example.com")
		mf := message.MaxForwards(70)
		contentLength := message.ContentLength(0)
		cseq := message.CSeq{Seq: 1, Method: message.REGISTER}

		rm := message.Request{
			Method:     message.REGISTER,
			RequestURI: reqURI.String(), //"sips:ss2.biloxi.example.com",
			Headers: []message.Header{
				&mf,
				&from,
				&to,
				&callid,
				&cseq,
				&contact,
				&contentLength,
			},
		}

		writer := OutputStub{}
		err := rm.Write(&writer)
		if err != nil {
			t.Errorf("failed writing to the output")
		}
		got := string(writer.value)

		if want != got {
			t.Errorf("request different from expected.\nwant: \n%s#############\ngot: \n%s#############", want, got)
		}
	})
}

func getRegisterMessageSample() string {
	var sb strings.Builder
	sb.WriteString("REGISTER sips:ss2.biloxi.example.com SIP/2.0\r\n")
	sb.WriteString("Max-Forwards: 70\r\n")
	sb.WriteString("From: Bob <sips:bob@biloxi.example.com>;tag=a73kszlfl\r\n")
	sb.WriteString("To: \"A. G. Bell\" <sips:agb@bell-telephone.com>\r\n")
	sb.WriteString("Call-ID: 1j9FpLxk3uxtm8tn@biloxi.example.com\r\n")
	sb.WriteString("CSeq: 1 REGISTER\r\n")
	sb.WriteString("Contact: Bob <sips:bob@biloxi.example.com>\r\n")
	sb.WriteString("Content-Length: 0\r\n")
	sb.WriteString("\r\n")
	return sb.String()
}
