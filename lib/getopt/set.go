package getopt

type argSet struct {
	expectedShort map[rune]*Option
	expectedLong  map[string]*Option
}

func newSet() *argSet {
	s := &argSet{
		expectedShort: make(map[rune]*Option),
		expectedLong:  make(map[string]*Option),
	}
	return s
}
