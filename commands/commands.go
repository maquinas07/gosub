package commands

import "github.com/maquinas07/gosub/lib/srt"

func Init() {
	init_time_shift()
	init_overlaps_join()
}

func Execute(subs *[]*srt.Subtitle) {
	perform_overlaps_join(&subs)
	perform_time_shifts(*subs)
}
