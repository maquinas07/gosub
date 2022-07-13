package getopt

import (
	"os"
)

type Args struct {
	program string
	options []*ValuedOption
	flags   []*Flag
}

var expectedArgs = newSet()

func (a *Args) GetOptions() []*ValuedOption {
	return a.options
}

func (a *Args) GetFlags() []*Flag {
	return a.flags
}

func (a *Args) GetProgram() string {
	return a.program
}

func (a *Args) parseShortOption(shortArg string) (isCompleteOption bool) {
	return false
}

func (a *Args) parseLongOption(longArg string) (isCompleteOption bool) {
	for _, v := range longArg {
		if v == '=' {
			// (*a.options)[longArg[:i]] = longArg[i+1:]
			return true
		}
	}
	return false
}

func (a *Args) addFlag(flag string) {
	// a.flags = append(a.flags, &flag)
}

// func parseValue(v interface{}, shortName rune, longName string) (opt Option) {
// 	// switch v.(type) {
// 	// case Option:
// 	// 	{
// 	// 		return &generic{v}
// 	// 	}
// 	// }
// }

func AddFlag(shortName rune, longName string, p *bool) (f Flag) {
	f = &flag{
		opt: &option{
			shortName: shortName,
			longName:  longName,
		},
		p: p,
	}
	var r Option = f
	expectedArgs.expectedShort[shortName] = &r
	expectedArgs.expectedLong[longName] = &r
	return
}

func AddOption(shortName rune, longName string, p *interface{}, isOptional bool, valueParser *ValueParser) (o Option) {
	o = &valuedOption{
		opt: &option{
			shortName: shortName,
			longName:  longName,
		},
		value:       p,
		valueParser: valueParser,
		isOptional:  isOptional,
	}
	var r Option = o
	expectedArgs.expectedShort[shortName] = &r
	expectedArgs.expectedLong[longName] = &r
	return
}

func Parse() (a *Args, err error) {
	a = &Args{
		program: os.Args[0],
	}

	var optOffset byte
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg[0] == '-' {
			if len(arg) == 1 {
				// Stdin arg
			}
			switch optOffset {
			case 0:
				{
					break
				}
			case 1, 2:
				{
					a.addFlag(os.Args[i-1][optOffset:])
					optOffset = 0
				}
			}
			if arg[1] == '-' {
				if !a.parseLongOption(arg[2:]) {
					optOffset = 2
				}
			} else {
				if !a.parseShortOption(arg[1:]) {
					optOffset = 1
				}
			}
		} else {
			switch optOffset {
			case 0:
				{
					//
				}
			case 1, 2:
				{
					// opt := &Option{
					// 	key:   os.Args[i-1][optOffset:],
					// 	value: arg,
					// }
					// a.options = append(a.options, opt)
				}
			}
		}
	}
	return
}
