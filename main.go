package main

import (
	"fmt"
	"os"

	"github.com/maquinas07/gosub/commands"
	"github.com/maquinas07/gosub/lib/getopt"
	"github.com/maquinas07/gosub/lib/srt"
)

func main() {
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
			return
		}
		defer inputFile.Close()
		subs, err := srt.ParseMemoryUnbound(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}

		commands.Execute(&subs)

		outputFile, err := os.Create("out." + params[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}
		defer outputFile.Close()
		err = srt.Save(subs, outputFile)
		if err != nil {
			fmt.Print("Error\n")
		}
	}
}
