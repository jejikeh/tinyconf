package main

import (
	"fmt"

	"github.com/jejikeh/ambient/common"
	"github.com/jejikeh/ambient/vm"
)

func main() {
	ambient := vm.NewAmbient()

	ambient.LoadProgram([]common.Instruction{
		common.NewPush(0),
		common.NewPush(1),
		common.NewDuplicate(1),
		common.NewDuplicate(1),
		common.NewPlus(),
		common.NewJump(2),
	})

	for i := 0; i < 100; i++ {
		err := ambient.Run()
		if err != common.Ok {
			fmt.Printf("Error: %s\n", err.String())
			ambient.PrintStack()
			panic(1)
		}
	}

	ambient.PrintStack()

	a, b := 0, 1
	for i := 0; i <= 64; i++ {
		fmt.Printf("	%d: %d\n", i, a)
		a, b = b, a+b
	}
	fmt.Println()
}