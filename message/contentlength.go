package message

import "strconv"

type ContentLength int

func (_ ContentLength) Tag() string {
	return "Content-Length"
}

func (cl ContentLength) String() string {
	return strconv.Itoa(int(cl))
}
