package utf8

var utf8bom = []byte{0xEF, 0xBB, 0xBF}

func StripUTF8BOM(data []byte) []byte {
	if len(data) < len(utf8bom) {
		return data
	}
	var i int
	for i = 0; i < len(utf8bom) && data[i] == utf8bom[i]; i++ {
	}
	return data[i:]
}
