package getopt

type argSet struct {
	expectedShort map[rune]interface{}
	expectedLong  map[string]interface{}
}

func newSet() *argSet {
	s := &argSet{
		expectedShort: make(map[rune]interface{}),
		expectedLong:  make(map[string]interface{}),
	}
	return s
}
