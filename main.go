package main

import (
	"flag"
	"fmt"

	"github.com/jejikeh/ambient/common"
	"github.com/jejikeh/ambient/vm"
)

func main() {
	sourcePath := flag.String("source", "", "Source file")
	output := flag.String("output", "", "Output file")

	flag.Parse()

	ambient := vm.NewAmbient()
	ambient.LoadAmbientAsmFromFile(*sourcePath)
	ambient.PrintInstructions()

	for !ambient.End {
		err := ambient.Run()
		if err != common.Ok {
			fmt.Printf("Error: %s\n", err.String())
			ambient.PrintStack()
			panic(1)
		}
	}

	ambient.PrintStack()
	ambient.SaveProgramToNewFile(*output)
}
