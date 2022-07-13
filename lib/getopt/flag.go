package getopt

type Flag interface {
	Option
	IsSet() bool
}

type flag struct {
	opt Option
	p   *bool
}

func (f *flag) set() {
	if f.p != nil {
		*f.p = true
	} else {
		v := true
		f.p = &v
	}
}

func (f *flag) IsSet() bool {
	if f.p == nil {
		return false
	}
	return *f.p
}

func (f *flag) ShortName() rune {
	return f.opt.ShortName()
}

func (f *flag) LongName() string {
	return f.opt.LongName()
}
