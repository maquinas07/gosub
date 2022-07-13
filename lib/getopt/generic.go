package getopt

import "errors"

type ValueParser func(s string) (*interface{}, error)

type ValuedOption interface {
	Option
	IsValueOptional() bool
	Value() *interface{}
}

type valuedOption struct {
	opt         Option
	value       *interface{}
	valueParser *ValueParser
	isOptional  bool
}

func (o *valuedOption) setValue(value string) (err error) {
	if o.valueParser != nil {
		(*o.value), err = (*o.valueParser)(value)
	} else {
		switch (*o.value).(type) {
		case *string:
			(*o.value) = &value
		default:
			err = errors.New("type is not supported")
		}
	}
	return
}

func (o *valuedOption) ShortName() rune {
	return o.opt.ShortName()
}

func (o *valuedOption) LongName() string {
	return o.opt.LongName()
}

func (o *valuedOption) Value() *interface{} {
	return o.value
}

func (o *valuedOption) IsValueOptional() bool {
	return o.isOptional
}
