package message

import "strconv"

type MaxForwards uint8

func (_ MaxForwards) Tag() string {
	return "Max-Forwards"
}

func (f MaxForwards) String() string {
	return strconv.FormatUint(uint64(f), 10)
}
