package getopt

type Option interface {
	ShortName() rune
	LongName() string
}

type option struct {
	shortName rune
	longName  string
}

func (o *option) ShortName() rune {
	return o.shortName
}

func (o *option) LongName() string {
	return o.longName
}
