package message

import "strconv"

type MaxForwards uint8

func (_ MaxForwards) Tag() string {
	return "Max-Forwards"
}

func (f MaxForwards) String() string {
	return strconv.FormatUint(uint64(f), 10)
}

func ParseMaxForwards(value string) (*MaxForwards, error) {
	v, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return nil, err
	}
	mf := MaxForwards(v)
	return &mf, nil
}
