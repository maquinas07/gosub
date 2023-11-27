package commands

import (
	"github.com/maquinas07/gosub/lib/getopt"
	"github.com/maquinas07/gosub/lib/srt"
)

var joinOverlaps bool

func initOverlapsJoin() {
	getopt.AddFlag('j', "join-overlaps", &joinOverlaps)
}

func performOverlapsJoin(subsReference **[]*srt.Subtitle) {
	if joinOverlaps {
		var newSubs []*srt.Subtitle
		subs := **subsReference
		for i := 0; i < len(subs); i++ {
			sub := &srt.Subtitle{
				Index:     subs[i].Index,
				StartTime: subs[i].StartTime,
				EndTime:   subs[i].EndTime,
				Dialogue:  subs[i].Dialogue,
			}
			for j := i + 1; j < len(subs) && subs[i].StartTime == subs[j].StartTime && subs[i].EndTime == subs[j].EndTime; j++ {
				sub.Dialogue = append(sub.Dialogue, subs[j].Dialogue...)
				i = j
			}
			newSubs = append(newSubs, sub)
			*subsReference = &newSubs
		}
	}
}
