package getopt

type flag struct {
	option
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

