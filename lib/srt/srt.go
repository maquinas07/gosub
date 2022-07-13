package srt

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/maquinas07/gosub/lib/ascii"
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

	{
		if scanner.Scan() {
			data := utf8.StripUTF8BOM(scanner.Bytes())
			p.parse(data)
		}
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
