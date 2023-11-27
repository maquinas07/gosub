package utf8

import "encoding/binary"

var utf8bom uint32 = 0xEFBBBF00

func StripUTF8BOM(data []byte) []byte {
	if utf8bom&binary.BigEndian.Uint32(data[:4]) == utf8bom {
		return data[3:]
	}
	return data
}
