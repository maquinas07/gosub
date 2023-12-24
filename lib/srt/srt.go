package srt

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/maquinas07/golibs/ascii"
	. "github.com/maquinas07/golibs/shared"
	"github.com/maquinas07/golibs/utf8"
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

var timingSeparator uint64 = ' '<<32 | '>'<<24 | '-'<<16 | '-'<<8 | ' '

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
	index, err = ascii.ParseInt(data)
	if err != nil {
		err = ErrInvalidIndex
		return
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
	if (binary.LittleEndian.Uint64(data[12:20]) & timingSeparator) != timingSeparator {
		err = ErrInvalidTiming
		return
	}
	startTime, err = parseTiming(data[:12])
	if err != nil {
		return
	}
	endTime, err = parseTiming(data[20:])
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
				if len(p.currentSubtitle.Dialogue) == 0 {
					fmt.Fprintf(os.Stderr, "warn: Found dialogue started by a newline character, ignoring\n")
					return true
				}
				p.subtitles = append(p.subtitles, p.currentSubtitle)
				p.currentState = newline
			}
		}
	}
	return true
}

func ParseMemoryUnbound(reader io.Reader) (subs []*Subtitle, err error) {
	fileBytes, err := io.ReadAll(reader)
	if err != nil {
		return
	}

	fileBytes = utf8.StripUTF8BOM(fileBytes)

	var p *parser = &parser{
		subtitles:    make([]*Subtitle, 0),
		currentState: newline,
	}

	var i, j int
	var gtg bool = true
	for ; i < len(fileBytes) && gtg; i++ {
		if fileBytes[i] == '\n' {
			gtg = p.parse(fileBytes[j:i])
			j = i + 1
		} else if fileBytes[i] == '\r' && i+1 < len(fileBytes) && fileBytes[i+1] == '\n' {
			gtg = p.parse(fileBytes[j:i])
			i = i + 1
			j = i + 1
		} else if len(fileBytes) == i+1 {
			p.parse(fileBytes[j : i+1])
			p.parse([]byte{})
		}
	}

	if p.err != nil {
		err = p.err
	}
	return p.subtitles, err
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

func fmtInt(buf []byte, v uint64) {
	w := len(buf) - 1
	if v == 0 {
		buf[w] = '0'
		return
	}
	for v > 0 && w >= 0 {
		buf[w] = byte(v%10) + '0'
		v /= 10
		w--
	}
}

// https://pvk.ca/Blog/2017/12/22/appnexus-common-framework-its-out-also-how-to-print-integers-faster/
func encodeTenThousands(hi, lo uint64) uint64 {
	merged := hi | (lo << 32)
	top := ((merged * 10486) >> 20) & ((0x7F << 32) | 0x7F)
	bot := merged - 100*top
	hundreds := (bot << 16) + top
	tens := (hundreds * 103) >> 10
	tens &= (0xF << 48) | (0xF << 32) | (0xF << 16) | 0xF
	tens += (hundreds - 10*tens) << 8

	return tens
}

func fmt16Uint64(buf []byte, v uint64) []byte {
	top := v / 100000000
	bottom := v % 100000000
	first :=
		0x3030303030303030 + encodeTenThousands(top/10000, top%10000)
	second :=
		0x3030303030303030 + encodeTenThousands(bottom/10000, bottom%10000)
	buf = binary.LittleEndian.AppendUint64(buf, first)
	buf = binary.LittleEndian.AppendUint64(buf, second)
	return buf
}

func serializeTimings(timing int64) (timings []byte) {
	hours := uint64(timing / Hour)
	minutes := uint64((timing % Hour) / Minute)
	seconds := uint64((timing % Minute) / Second)
	millis := uint64((timing % Second) / Millisecond)

	timings = make([]byte, 0, 16)
	timings = fmt16Uint64(timings, hours*10000000000+minutes*10000000+seconds*10000+millis)[4:]

	timings[8] = ','
	timings[5] = ':'
	timings[2] = ':'
	return timings
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
		c = append(c, byte(timingSeparator),
			byte(timingSeparator>>8),
			byte(timingSeparator>>16),
			byte(timingSeparator>>24),
			byte(timingSeparator>>32))
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
