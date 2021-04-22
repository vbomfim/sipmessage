package message

import "strconv"

//ContentLength represents the Content-Length SIP Header
type ContentLength int

//Header returns the header tag
func (ContentLength) Header() string {
	return "Content-Length"
}

//CHeader returns the compact header tag
func (ContentLength) CHeader() string {
	return "l"
}

func (cl ContentLength) String() string {
	return strconv.Itoa(int(cl))
}

//ParseContentLength parses a string value to a ContentLength instance and returns as a Header interface
func ParseContentLength(value string) (*ContentLength, error) {
	v, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	cl := ContentLength(v)
	return &cl, nil
}
