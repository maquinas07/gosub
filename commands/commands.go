package commands

import (
	"github.com/maquinas07/gosub/lib/srt"
)

func Init() {
	initTimeShift()
	initOverlapsJoin()
}

func Execute(subs *[]*srt.Subtitle) {
	performOverlapsJoin(&subs)
	performTimeShifts(*subs)
}
