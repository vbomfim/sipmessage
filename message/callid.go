package message

type CallID string

func (_ CallID) Tag() string {
	return "Call-ID"
}

func (c CallID) String() string {
	return string(c)
}
