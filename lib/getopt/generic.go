package getopt

import "errors"

type ValueParser func(s string) (interface{}, error)

type valuedOption struct {
	option
	value       interface{}
	valueParser ValueParser
	isOptional  bool
}

func (o *valuedOption) setValue(value string) (err error) {
	if o.valueParser != nil {
		if o.value != nil {
			o.value, err = o.valueParser(value)
		} else {
			_, err = o.valueParser(value)
		}
	} else {
		switch o.value.(type) {
		case *string:
			o.value = &value
		default:
			err = errors.New("type is not supported")
		}
	}
	return
}
