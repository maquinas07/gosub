package srt

import (
	"bufio"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var subs []*Subtitle
var timings []byte

func BenchmarkFmt1(b *testing.B) {
    for i := 0; i < b.N; i++ {
		timings = serializeTimings(rand.Int63n(3 * 1000000000))
	}
}

func BenchmarkSrt(b *testing.B) {
	b.StopTimer()
	fd, err := os.Open("./test.srt")
	if err != nil {
		b.FailNow()
	}
	defer fd.Close()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		subs, err = Parse(fd)
	}
}

func BenchmarkSrtUnbound(b *testing.B) {
	b.StopTimer()
	fd, err := os.Open("./test.srt")
	if err != nil {
		b.FailNow()
	}
	defer fd.Close()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		subs, err = ParseMemoryUnbound(fd)
	}
}

func TestSrt(t *testing.T) {
	fd, err := os.Open("./test.srt")
	if err != nil {
		t.Fatalf("Couldn't open test file")
	}

	subs, err := ParseMemoryUnbound(fd)
	if err != nil {
		t.Fatalf("Fatal error in parsing %s", err)
	}
	assert.NotEmpty(t, subs)

	of, err := os.Create("./output.srt")
	if err != nil {
		t.Fatalf("Couldn't create output file")
	}
	Save(subs, of)
	fd.Close()
	of.Close()

	fd, err = os.Open("./test.srt")
	if err != nil {
		t.Fatalf("Couldn't open test file")
	}
	of, err = os.Open("./output.srt")
	if err != nil {
		t.Fatalf("Couldn't open output test file for reading")
	}

	fdReader := bufio.NewScanner(fd)
	ofReader := bufio.NewScanner(of)
	for fdReader.Scan() && ofReader.Scan() {
		assert.Nil(t, fdReader.Err())
		assert.Nil(t, ofReader.Err())
		assert.Subset(t, fdReader.Bytes(), ofReader.Bytes())
	}

	fdReader.Scan()
	ofReader.Scan()
	assert.Empty(t, fdReader.Bytes())
	assert.Empty(t, ofReader.Bytes())
	assert.Nil(t, fdReader.Err())
	assert.Nil(t, ofReader.Err())
}

func FuzzTimingsParser(f *testing.F) {
    f.Add([]byte("00:00:13,596 --> 00:00:16,141"))
    f.Fuzz(func (t *testing.T, timing []byte) {
        startTime, endTime, err := parseTimings(timing)
        if (startTime < 0 || endTime < 0) && err == nil {
            t.Errorf("Something fishy happened: %v, %v\n", subs, err)
        }
    })
}

func FuzzTimingsSerializer(f *testing.F) {
    f.Fuzz(func (t *testing.T, timing int64) {
        serializeTimings(timing)
    })
}
