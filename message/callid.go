package message

type CallID string

func (_ CallID) Tag() string {
	return "Call-ID"
}

func (_ CallID) CTag() string {
	return "i"
}

func (c CallID) String() string {
	return string(c)
}

func ParseCallID(value string) (*CallID, error) {
	c := CallID(value)
	return &c, nil
}
