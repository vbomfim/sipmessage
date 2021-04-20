package message

import "strconv"

type ContentLength int

func (_ ContentLength) Tag() string {
	return "Content-Length"
}

func (_ ContentLength) CTag() string {
	return "l"
}

func (cl ContentLength) String() string {
	return strconv.Itoa(int(cl))
}

func ParseContentLength(value string) (*ContentLength, error) {
	v, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	cl := ContentLength(v)
	return &cl, nil
}
