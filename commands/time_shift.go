package commands

import (
	"errors"

	"github.com/maquinas07/gosub/lib/ascii"
	"github.com/maquinas07/gosub/lib/getopt"
	. "github.com/maquinas07/gosub/lib/shared"
	"github.com/maquinas07/gosub/lib/srt"
)

const (
	invalidTimeFormatErrorMessage = "Invalid time format."
)

type TimeSegment struct {
	StartTime int64
	EndTime   int64
}

type TimeShift struct {
	ShiftByMs   int64
	TimeFilters []TimeSegment
}

func (o *TimeShift) shift(subs []*srt.Subtitle) {
	var currentTimeFilter *TimeSegment
	var j int = 0
	if len(o.TimeFilters) > j {
		currentTimeFilter = &o.TimeFilters[j]
		j++
	}
	for i := 0; i < len(subs); i++ {
		sub := subs[i]
		if currentTimeFilter == nil || currentTimeFilter.StartTime > sub.StartTime && currentTimeFilter.EndTime < sub.EndTime {
			sub.StartTime += o.ShiftByMs
			sub.EndTime += o.ShiftByMs
			if sub.StartTime < 0 {
				sub.StartTime = 0
			}
			if sub.EndTime < 0 {
				sub.EndTime = 0
			}
			if len(o.TimeFilters) > j {
				currentTimeFilter = &o.TimeFilters[j]
				j++
			}
		}
	}
}

var timeShifts []TimeShift

func initTimeShift() {
	getopt.AddOption('s', "shift", nil, false, func(s string) (dummy interface{}, err error) {
		var parsedTimeShift TimeShift
		var exp int64 = 0
		var base int64 = 1
		for i := len(s) - 1; i >= 0; i-- {
			switch s[i] {
			case 's':
				fallthrough
			case 'S':
				{
					if exp > 0 && base == 1 {
						err = errors.New(invalidTimeFormatErrorMessage)
						return
					}
					if base > 1 {
						parsedTimeShift.ShiftByMs *= exp
					}
					base = 1
					exp = Second
					break
				}
			case 'm':
				fallthrough
			case 'M':
				{
					if base > 1 {
						parsedTimeShift.ShiftByMs *= exp
						exp = 0
					}
					base = 1
					if exp == 1000 {
						exp = Millisecond
					} else if exp == 0 {
						exp = Minute
					} else {
						err = errors.New(invalidTimeFormatErrorMessage)
						return
					}
					break
				}
			case 'h':
				fallthrough
			case 'H':
				{
					if exp > 0 && base == 1 {
						err = errors.New(invalidTimeFormatErrorMessage)
						return
					}
					if base > 1 {
						parsedTimeShift.ShiftByMs *= exp
					}
					base = 1
					exp = Hour
					break
				}
			case '-':
				{
					if i == 0 {
						parsedTimeShift.ShiftByMs *= -1
					} else {
						err = errors.New(invalidTimeFormatErrorMessage)
						return
					}
					break
				}
			case '+':
				{
					if i != 0 {
						err = errors.New(invalidTimeFormatErrorMessage)
					}
					break
				}
			default:
				{
					if exp == 0 {
						err = errors.New(invalidTimeFormatErrorMessage)
						return
					}
					var value int
					value, err = ascii.ToDigit(s[i])
					if err != nil {
						return
					}
					parsedTimeShift.ShiftByMs += base * int64(value)
					base *= 10
					break
				}
			}
		}
		parsedTimeShift.ShiftByMs *= exp
		timeShifts = append(timeShifts, parsedTimeShift)
		return nil, err
	})
}

func performTimeShifts(subs []*srt.Subtitle) {
	for i := 0; i < len(timeShifts); i++ {
		timeShifts[i].shift(subs)
	}
}
