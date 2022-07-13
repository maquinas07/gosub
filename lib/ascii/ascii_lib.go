package ascii

import "fmt"

func isWhitespace(char byte) bool {
	return char == '\n' || char == ' ' || char == '\t'
}

func TrimWhitespaces(data []byte) []byte {
	var i, j int
	for i = len(data) - 1; isWhitespace(data[i]); i-- {
	}
	for j = 0; isWhitespace(data[j]); j++ {
	}
	return data[j : i+1]
}

func IsDigit(data byte) bool {
	return data >= 48 && data <= 57
}

func ToDigit(data byte) (digit int, err error) {
	if !IsDigit(data) {
		err = fmt.Errorf("ascii_lib: invalid digit")
		return
	}
	digit = int(data - 48)
	return
}
