package main

import (
	"fmt"
	"os"

	"github.com/maquinas07/gosub/lib/getopt"
	"github.com/maquinas07/gosub/lib/srt"
)

func main() {
	args := os.Args
	getopt.Parse()
	fs, err := os.Open("./test.srt")
	if err != nil {
		return
	}
	subs, err := srt.Parse(fs)
	if err != nil {
		return
	}
	fmt.Print(subs, args)
}
