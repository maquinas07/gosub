package main

import (
	"fmt"
	"os"

	"github.com/maquinas07/gosub/lib/getopt"
	"github.com/maquinas07/gosub/lib/srt"
)

func main() {
	_, err := getopt.Parse()
	if err != nil {
		return
	}

	inputFile, err := os.Open("./test.srt")
	if err != nil {
		return
	}
	defer inputFile.Close()
	subs, err := srt.Parse(inputFile)
	if err != nil {
		return
	}
	var newSubs []*srt.Subtitle
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
	}

	outputFile, err := os.Create("./output.srt")
	if err != nil {
		return
	}
	defer outputFile.Close()
	err = srt.Save(newSubs, outputFile)
	if err != nil {
		fmt.Print("Error\n")
	}
}
