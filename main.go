package main

import (
	"fmt"
	"os"
	"path"

	"github.com/maquinas07/golibs/getopt"
	"github.com/maquinas07/gosub/commands"
	"github.com/maquinas07/gosub/lib/srt"
)

var replace bool

func addCommonFlags() {
	getopt.AddFlag('r', "replace", &replace)
}

func main() {
	addCommonFlags()
	commands.Init()

	args, err := getopt.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	params := args.GetParams()
	for i := 0; i < len(params); i++ {
		inputFile, err := os.Open(params[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}
		defer inputFile.Close()
		subs, err := srt.ParseMemoryUnbound(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}

		commands.Execute(&subs)

		var outputFile *os.File
		defer outputFile.Close()
		if !replace {
			outputFile, err = os.Create("out." + path.Base(params[i]))
		} else {
			outputFile, err = os.Create(params[i])
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}
		err = srt.Save(subs, outputFile)
		if err != nil {
			fmt.Print("Error\n")
		}
	}
}
