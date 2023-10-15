package main

import (
	"fmt"

	"github.com/jejikeh/ambient/common"
	"github.com/jejikeh/ambient/vm"
)

func main() {
	x := `psh 0
psh 1
dplc 1
dplc 1
plus
jmp 2`

	ambient := vm.NewAmbient()
	ambient.LoadAmbientAsm(x)
	ambient.PrintInstructions()

	for i := 0; i < 100; i++ {
		err := ambient.Run()
		if err != common.Ok {
			fmt.Printf("Error: %s\n", err.String())
			ambient.PrintStack()
			panic(1)
		}
	}

	ambient.PrintStack()
	ambient.SaveProgramToNewFile("examples/binary/", "fib.amb")
}
