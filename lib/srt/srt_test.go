package srt

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var subs []*Subtitle

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

func TestSrt(t *testing.T) {
	fd, err := os.Open("./test.srt")
	if err != nil {
		t.FailNow()
	}

	subs, err := Parse(fd)
	if err != nil {
		t.FailNow()
	}
	assert.NotEmpty(t, subs)

	of, err := os.Create("./output.srt")
	if err != nil {
		t.FailNow()
	}
	Save(subs, of)
	fd.Close()
	of.Close()

	fd, err = os.Open("./test.srt")
	if err != nil {
		t.FailNow()
	}
	of, err = os.Open("./output.srt")
	if err != nil {
		t.FailNow()
	}

	fdReader := bufio.NewScanner(fd)
	ofReader := bufio.NewScanner(of)
	for fdReader.Scan() && ofReader.Scan() {
		assert.Nil(t, fdReader.Err())
		assert.Nil(t, ofReader.Err())
		assert.Equal(t, fdReader.Bytes(), ofReader.Bytes())
	}

	fdReader.Scan()
	ofReader.Scan()
	assert.Empty(t, fdReader.Bytes())
	assert.Empty(t, ofReader.Bytes())
	assert.Nil(t, fdReader.Err())
	assert.Nil(t, ofReader.Err())
}
