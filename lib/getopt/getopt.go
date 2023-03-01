package getopt

import (
	"fmt"
	"os"
	"strings"
)

type Args struct {
	program string
	params  []string
}

var expectedArgs = newSet()

func (a *Args) GetProgram() string {
	return a.program
}

func (a *Args) GetParams() []string {
	return a.params
}

func AddFlag(shortName rune, longName string, p *bool) {
    f := &flag{
        option: option{
			shortName: shortName,
			longName:  longName,
		},
		p: p,
	}
	expectedArgs.expectedShort[shortName] = f
	expectedArgs.expectedLong[longName] = f
	return
}

func AddOption(shortName rune, longName string, p interface{}, isOptional bool, valueParser ValueParser) {
    o := &valuedOption{
        option: option {
			shortName: shortName,
			longName:  longName,
		},
		value:       p,
		valueParser: valueParser,
		isOptional:  isOptional,
	}
	expectedArgs.expectedShort[shortName] = o
	expectedArgs.expectedLong[longName] = o
	return
}

func Parse() (a *Args, err error) {
	args := os.Args
	a = &Args{
		program: args[0],
	}

	for i := 1; i < len(args); i++ {
		arg := args[i]
		if arg[0] == '-' {
			if arg[1] == '-' {
				i := strings.IndexRune(arg, '=')
				var value string
				if i > 0 {
					arg = arg[:i]
					value = arg[i+1:]
				}
				opt := expectedArgs.expectedLong[arg]

				if opt != nil {
					valueOpt, ok := opt.(*valuedOption)
					if ok {
						if value != "" {
							err = valueOpt.setValue(value)
                            if err != nil {
                                return
                            }
						} else if !valueOpt.isOptional && i < 0 {
							if len(args) < i+2 {
								err = fmt.Errorf("missing argument value")
								return
							}
							i++
							valueOpt.setValue(args[i])
							continue
						}
					} else {
						flag, ok := opt.(*flag)
						if !ok {
							a.params = append(a.params, arg)
							continue
						}
						flag.set()
					}
				}
			} else {
				for j, r := range arg {
					opt := expectedArgs.expectedShort[r]
					if opt != nil {
						flag, ok := opt.(*flag)
						if ok {
							flag.set()
						} else {
							valueOpt, ok := opt.(*valuedOption)
							if !ok {
								a.params = append(a.params, arg)
								continue
							}
							value := arg[j+1:]
							if value == "" && !valueOpt.isOptional {
								if len(args) < i+2 {
									err = fmt.Errorf("missing argument value")
									return
								}
								i++
								value = args[i]
							}
							err = valueOpt.setValue(value)
                            if err != nil {
                                return
                            }
						}
					}
				}
			}
		} else {
			a.params = append(a.params, arg)
		}
	}
	return
}
