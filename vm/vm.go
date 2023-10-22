package vm

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/jejikeh/ambient/common"
)

func (a *VirtualMachine) LoadProgram(program []Instruction) {
	a.Instructions = program
}

func (a *VirtualMachine) Run() common.Error {
	if a.InstructionPointer < 0 || a.InstructionPointer >= len(a.Instructions) {
		return common.IllegalInstruction
	}

	instruction := a.Instructions[a.InstructionPointer]

	switch instruction.Type {
	case Push:
		// Push a value onto the stack.
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. PSH 1
		//		2. PRINT_STACK: [0, 1, 1]

		a.Stack = append(a.Stack, instruction.Operand)
		a.InstructionPointer++

	case Duplicate:
		// Duplicate the top of the stack.
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. DPLC 0
		//		2. PRINT_STACK: [0, 1, 0]

		if len(a.Stack)-instruction.Operand <= 0 {
			return common.StackUnderflow
		}

		if instruction.Operand < 0 {
			return common.IllegalInstruction
		}

		a.Stack = append(a.Stack, a.Stack[len(a.Stack)-1-instruction.Operand])
		a.InstructionPointer++

	case Plus:
		// Add the top two values on the stack.
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. PSH 1
		// 		2. PSH 1
		// 		3. ADD
		//		4. PRINT_STACK: [0, 1, 2]

		if len(a.Stack) < 2 {
			return common.StackUnderflow
		}

		a.Stack[len(a.Stack)-2] = a.Stack[len(a.Stack)-2] + a.Stack[len(a.Stack)-1]
		a.Stack = a.Stack[:len(a.Stack)-1]
		a.InstructionPointer++

	case Minus:
		// Subtract the top two values on the stack.
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. PSH 2
		// 		2. PSH 1
		// 		3. SUB
		// 		4. PRINT_STACK: [0, 1, 1]

		if len(a.Stack) < 2 {
			return common.StackUnderflow
		}

		a.Stack[len(a.Stack)-2] = a.Stack[len(a.Stack)-2] - a.Stack[len(a.Stack)-1]
		a.Stack = a.Stack[:len(a.Stack)-1]
		a.InstructionPointer++

	case Multiply:
		// Multiply the top two values on the stack.
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. PSH 2
		// 		2. PSH 2
		// 		3. MUL
		// 		4. PRINT_STACK: [0, 1, 4]

		if len(a.Stack) < 2 {
			return common.StackUnderflow
		}

		a.Stack[len(a.Stack)-2] = a.Stack[len(a.Stack)-2] * a.Stack[len(a.Stack)-1]
		a.Stack = a.Stack[:len(a.Stack)-1]
		a.InstructionPointer++

	case Divide:
		// Divide the top two values on the stack.
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. PSH 4
		// 		2. PSH 2
		// 		3. DIV
		// 		4. PRINT_STACK: [0, 1, 2]

		if len(a.Stack) < 2 {
			return common.StackUnderflow
		}

		if a.Stack[len(a.Stack)-1] == 0 {
			return common.DivisionByZero
		}

		a.Stack[len(a.Stack)-2] = a.Stack[len(a.Stack)-2] / a.Stack[len(a.Stack)-1]
		a.Stack = a.Stack[:len(a.Stack)-1]
		a.InstructionPointer++

	case Jump:
		// Jump to a new instruction.
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. JMP 2
		// 		2. PRINT_STACK: [0, 1]
		if instruction.Operand < 0 || instruction.Operand >= len(a.Instructions) {
			return common.IllegalInstruction
		}

		a.InstructionPointer = instruction.Operand

	case JumpIfTrue:
		// Jump to a new instruction if the top of the stack is true (1).
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. PSH 0
		// 		2. PSH 1
		// 		3. JMPIF 4
		// 		4. PRINT_STACK: [0, 1]

		if len(a.Stack) < 1 {
			return common.StackUnderflow
		}

		if a.Stack[len(a.Stack)-1] != 1 {
			a.InstructionPointer++
			break
		}

		a.Stack = a.Stack[:len(a.Stack)-1]
		a.InstructionPointer = instruction.Operand

	case Equal:
		// instruction if the top of the stack is equal.
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. PSH 1
		// 		2. PSH 1
		// 		3. EQ
		// 		4. PRINT_STACK: [0, 1, 1]
		if len(a.Stack) < 2 {
			return common.StackUnderflow
		}

		eq := (a.Stack[len(a.Stack)-2] == a.Stack[len(a.Stack)-1])
		if eq {
			a.Stack[len(a.Stack)-2] = 1
		} else {
			a.Stack[len(a.Stack)-2] = 0
		}

		a.Stack = a.Stack[:len(a.Stack)-1]
		a.InstructionPointer++

	case End:
		a.InstructionPointer++

	default:
		return common.IllegalInstruction
	}

	return common.Ok
}

func (a *VirtualMachine) Execute(executingLimit int, printCurrentInstruction bool) {
	isInfinite := executingLimit < 0
	for i := 0; (i < executingLimit && (a.Instructions[a.InstructionPointer].Type != End)) || isInfinite; i++ {
		err := a.Run()
		if err != common.Ok {
			color.Set(color.FgHiRed)
			defer color.Unset()

			log.Printf("Error: %s\n", err.String())
			a.PrintStack()
			panic(1)
		}

		if printCurrentInstruction {
			fmt.Printf("[%d] Current pointer -> [%d: %s]\n", i, a.InstructionPointer, a.Instructions[a.InstructionPointer].Type.String())
		}
	}
}

func (a *VirtualMachine) PrintStack() {
	fmt.Println("Stack:")
	if len(a.Stack) == 0 {
		fmt.Println("	Stack is empty")
	}

	for i, v := range a.Stack {
		fmt.Printf("	%d: %d\n", i, v)
	}

	fmt.Println()
}

func (a *VirtualMachine) PrintInstructions() {
	fmt.Println("Instructions:")
	for i, v := range a.Instructions {
		fmt.Printf("	%d: %s %d\n", i, v.Type.String(), v.Operand)
	}

	fmt.Println()
}
