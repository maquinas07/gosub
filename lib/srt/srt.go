package srt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/maquinas07/gosub/lib/ascii"
	. "github.com/maquinas07/gosub/lib/shared"
	"github.com/maquinas07/gosub/lib/utf8"
)

// The parser is a state machine
type parserState uint

const (
	newline  parserState = iota // followed by NEWLINE or INDEX
	index                       // followed by DIALOGUE
	dialogue                    // followed by NEWLINE
)

type parser struct {
	currentState    parserState
	currentSubtitle *Subtitle
	subtitles       []*Subtitle
	err             error
}

var timingSeparator = []byte{' ', '-', '-', '>', ' '}

type Subtitle struct {
	Index     int
	StartTime int64
	EndTime   int64
	Dialogue  []byte
}

var (
	ErrInvalidIndex  = errors.New("gosub srt: invalid index in srt subtitle")
	ErrInvalidTiming = errors.New("gosub srt: invalid timing in srt subtitle")
)

func parseIndex(data []byte) (index int, err error) {
	for i, baseExp := len(data)-1, 1; i >= 0; i, baseExp = i-1, baseExp*10 {
		var value int
		value, err = ascii.ToDigit(data[i])
		if err != nil {
			err = ErrInvalidIndex
			return
		}
		index += int(value) * baseExp
	}
	return
}

func parseTiming(data []byte) (time int64, err error) {
	baseExp := 1
	timeMult := 1
	charCount := 0
	for i := len(data) - 1; i >= 0; i-- {
		switch data[i] {
		case ':':
			if charCount != 2 {
				err = ErrInvalidTiming
				return
			}
			charCount = 0
			timeMult *= 6
			baseExp /= 10
		case ',':
			if charCount != 3 {
				err = ErrInvalidTiming
				return
			}
			charCount = 0
		default:
			var value int
			value, err = ascii.ToDigit(data[i])
			if err != nil {
				err = ErrInvalidTiming
				return
			}
			time += int64(value * baseExp * timeMult)
			baseExp *= 10
			charCount++
		}
	}
	return
}

func parseTimings(data []byte) (startTime int64, endTime int64, err error) {
	if len(data) != 29 {
		err = ErrInvalidTiming
		return
	}
	var i, j int
	for i, j = 12, 0; j < len(timingSeparator) && data[i] == timingSeparator[j]; i, j = i+1, j+1 {
	}
	if j != len(timingSeparator) {
		err = ErrInvalidTiming
		return
	}
	startTime, err = parseTiming(data[:12])
	if err != nil {
		return
	}
	endTime, err = parseTiming(data[i:])
	return
}

func (p *parser) parse(data []byte) bool {
	if p.err != nil {
		return false
	}
	switch p.currentState {
	case newline:
		{
			if len(data) > 0 {
				p.currentSubtitle = new(Subtitle)
				p.currentSubtitle.Index, p.err = parseIndex(data)
				if p.err != nil {
					return false
				}
				p.currentState = index
			}
		}
	case index:
		{
			p.currentSubtitle.StartTime, p.currentSubtitle.EndTime, p.err = parseTimings(data)
			if p.err != nil {
				return false
			}
			p.currentState = dialogue
		}
	case dialogue:
		{
			if len(data) > 0 {
				p.currentSubtitle.Dialogue = append(p.currentSubtitle.Dialogue, data...)
				p.currentSubtitle.Dialogue = append(p.currentSubtitle.Dialogue, '\n')
			} else {
				p.subtitles = append(p.subtitles, p.currentSubtitle)
				p.currentState = newline
			}
		}
	}
	return true
}

func Parse(reader io.Reader) (subs []*Subtitle, err error) {
	scanner := bufio.NewScanner(reader)

	var line int
	var p *parser = &parser{
		subtitles:    make([]*Subtitle, 0),
		currentState: newline,
	}

	if scanner.Scan() {
		data := utf8.StripUTF8BOM(scanner.Bytes())
		p.parse(data)
	}

	for scanner.Scan() && p.parse(scanner.Bytes()) {
		line++
	}

	if scanner.Err() != nil {
		err = fmt.Errorf("gosub: failed to read from reader %d. error in line %d", reader, line)
		return
	}
	if p.err != nil {
		err = p.err
	}
	subs = p.subtitles
	return
}

func fmtInt(buf []byte, v uint64) int {
	w := len(buf) - 1
	if v == 0 {
		buf[w] = '0'
	} else {
		for v > 0 {
			buf[w] = byte(v%10) + '0'
			v /= 10
			w--
		}
	}
	return w
}

func serializeTimings(timing int64) (timings []byte) {
	hours := uint64(timing / Hour)
	minutes := uint64((timing % Hour) / Minute)
	seconds := uint64((timing % Minute) / Second)
	millis := uint64((timing % Second) / Millisecond)
	timings = make([]byte, 12)
	for i := 0; i < len(timings); i++ {
		timings[i] = '0'
	}
	fmtInt(timings, millis)
	fmtInt(timings[:8], seconds)
	fmtInt(timings[:5], minutes)
	fmtInt(timings[:2], hours)
	timings[8] = ','
	timings[5] = ':'
	timings[2] = ':'
	return
}

func Save(subs []*Subtitle, writer io.Writer) (err error) {
	w := bufio.NewWriter(writer)
	for i, v := range subs {
		var c []byte
		index := make([]byte, 1+int(math.Log10(float64(i+1))))
		fmtInt(index, uint64(i+1))
		c = append(c, index...)
		c = append(c, '\n')
		c = append(c, serializeTimings(v.StartTime)...)
		c = append(c, timingSeparator...)
		c = append(c, serializeTimings(v.EndTime)...)
		c = append(c, '\n')
		c = append(c, v.Dialogue...)
		c = append(c, '\n')
		_, err = w.Write(c)
		if err != nil {
			return
		}
	}

	err = w.Flush()
	return
}
