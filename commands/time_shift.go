package commands

import (
	"errors"

	"github.com/maquinas07/golibs/ascii"
	"github.com/maquinas07/golibs/getopt"
	. "github.com/maquinas07/golibs/shared"
	"github.com/maquinas07/gosub/lib/srt"
)

var (
	errInvalidTimeFormat = errors.New("Invalid time format.\n")
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
	for i := 0; i < len(subs); i++ {
		sub := subs[i]
		var currentTimeFilter *TimeSegment
		for j := 0; j <= len(o.TimeFilters); j++ {
			if len(o.TimeFilters) > j {
				currentTimeFilter = &o.TimeFilters[j]
			}
			if currentTimeFilter == nil || (currentTimeFilter.StartTime < 0 || sub.StartTime >= currentTimeFilter.StartTime) && (currentTimeFilter.EndTime < 0 || sub.EndTime <= currentTimeFilter.EndTime) {
				sub.StartTime += o.ShiftByMs
				sub.EndTime += o.ShiftByMs
				if sub.StartTime < 0 {
					sub.StartTime = 0
				}
				if sub.EndTime < 0 {
					sub.EndTime = 0
				}
				break
			}
		}
	}
}

var timeShifts []TimeShift

func initTimeShift() {
	getopt.AddOption('s', "shift", nil, false, func(s string) (dummy interface{}, err error) {
		var parsedTimeShift TimeShift
		var currentTimeFilter *TimeSegment
		var accumulatedTime int64
		var exp int64 = 0
		var base int64 = 1
		for i := len(s) - 1; i >= 0; i-- {
			switch s[i] {
			case 's':
				fallthrough
			case 'S':
				{
					if exp > 0 && base == 1 {
						err = errInvalidTimeFormat
						return
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
						exp = 0
					}
					base = 1
					if exp == 1000 {
						exp = Millisecond
					} else if exp == 0 {
						exp = Minute
					} else {
						err = errInvalidTimeFormat
						return
					}
					break
				}
			case 'h':
				fallthrough
			case 'H':
				{
					if exp > 0 && base == 1 {
						err = errInvalidTimeFormat
						return
					}
					if base > 1 {
						accumulatedTime *= exp
					}
					base = 1
					exp = Hour
					break
				}
			case '-':
				{
					if i == 0 {
						accumulatedTime *= -1
					} else {
						err = errInvalidTimeFormat
						return
					}
					break
				}
			case '+':
				{
					if i != 0 {
						err = errInvalidTimeFormat
					}
					break
				}
			case '>':
				{
					if currentTimeFilter != nil && currentTimeFilter.StartTime > 0 {
						parsedTimeShift.TimeFilters = append(parsedTimeShift.TimeFilters, *currentTimeFilter)
					}
					currentTimeFilter = new(TimeSegment)
					currentTimeFilter.EndTime = -1
					currentTimeFilter.StartTime = accumulatedTime
					accumulatedTime = 0
					exp = 0
					base = 1
					break
				}
			case '<':
				{
					if currentTimeFilter != nil && currentTimeFilter.EndTime > 0 {
						parsedTimeShift.TimeFilters = append(parsedTimeShift.TimeFilters, *currentTimeFilter)
					}
					currentTimeFilter = new(TimeSegment)
					currentTimeFilter.StartTime = -1
					currentTimeFilter.EndTime = accumulatedTime
					accumulatedTime = 0
					exp = 0
					base = 1
					break
				}
			default:
				{
					if exp == 0 {
						err = errInvalidTimeFormat
						return
					}
					var value int
					value, err = ascii.ToDigit(s[i])
					if err != nil {
						return
					}
					accumulatedTime += base * int64(value) * exp
					base *= 10
					break
				}
			}
		}
		if currentTimeFilter != nil {
			parsedTimeShift.TimeFilters = append(parsedTimeShift.TimeFilters, *currentTimeFilter)
		}
		parsedTimeShift.ShiftByMs = accumulatedTime
		timeShifts = append(timeShifts, parsedTimeShift)
		return nil, err
	})
}

func performTimeShifts(subs []*srt.Subtitle) {
	for i := 0; i < len(timeShifts); i++ {
		timeShifts[i].shift(subs)
	}
}
