package message

import "strconv"

//MaxForwards represents the Max-Forwards SIP Header
type MaxForwards uint8

//Header returns the header tag
func (MaxForwards) Header() string {
	return "Max-Forwards"
}

//CHeader returns the compact header tag
func (MaxForwards) CHeader() string {
	return "Max-Forwards"
}

func (f MaxForwards) String() string {
	return strconv.FormatUint(uint64(f), 10)
}

//ParseMaxForwards a string value to a MaxForwards instance
func ParseMaxForwards(value string) (*MaxForwards, error) {
	v, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return nil, err
	}
	mf := MaxForwards(v)
	return &mf, nil
}
