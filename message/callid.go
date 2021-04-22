package message

//CallID represents the Call-ID SIP header
type CallID string

//Header returns the header tag
func (CallID) Header() string {
	return "Call-ID"
}

//CHeader returns the compact header tag
func (CallID) CHeader() string {
	return "i"
}

//String returns the CallID in string format
func (c CallID) String() string {
	return string(c)
}

//ParseCallID receives a string and parses it to the CallID type
func ParseCallID(value string) (*CallID, error) {
	c := CallID(value)
	return &c, nil
}
