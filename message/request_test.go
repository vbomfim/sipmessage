package message

import (
	"strings"
	"testing"
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

		rm := Request{
			Method:     REGISTER,
			RequestURI: "sips:ss2.biloxi.example.com",
			Headers: []HeaderField{
				{MaxForwards, []rune("70")},
				{From, []rune("Bob <sips:bob@biloxi.example.com>;tag=a73kszlfl")},
				{To, []rune("Bob <sips:bob@biloxi.example.com>")},
				{CallID, []rune("1j9FpLxk3uxtm8tn@biloxi.example.com")},
				{CSeq, []rune("1 REGISTER")},
				{Contact, []rune("<sips:bob@client.biloxi.example.com>")},
				{ContentLength, []rune("0")},
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
	sb.WriteString("MAX-FORWARDS: 70\r\n")
	sb.WriteString("FROM: Bob <sips:bob@biloxi.example.com>;tag=a73kszlfl\r\n")
	sb.WriteString("TO: Bob <sips:bob@biloxi.example.com>\r\n")
	sb.WriteString("CALL-ID: 1j9FpLxk3uxtm8tn@biloxi.example.com\r\n")
	sb.WriteString("CSEQ: 1 REGISTER\r\n")
	sb.WriteString("CONTACT: <sips:bob@client.biloxi.example.com>\r\n")
	sb.WriteString("CONTENT-LENGTH: 0\r\n")
	sb.WriteString("\r\n")
	return sb.String()
}
