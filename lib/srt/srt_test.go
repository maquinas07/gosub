package srt

import (
	"fmt"
	"os"
	"testing"
)

func TestSrt(t *testing.T) {
	fd, err := os.Open("/media/Ricardazo/Anime/86/[Erai-raws] 86 (2021) - 01 [1080p][Multiple Subtitle][E2CC3D55].ja.srt")
	if err != nil {
		t.Fail()
	}
	subs, err := Parse(fd)
	if err != nil {
		t.Fail()
	}
	for _, v := range subs {
		if v == nil {
			t.Fail()
		} else {
			fmt.Printf("%d\n", v.StartTime)
		}
	}
}
