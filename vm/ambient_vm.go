package vm

import (
	"fmt"

	"github.com/jejikeh/ambient/common"
)

type Ambient struct {
	Stack              []int
	Instructions       []common.Instruction
	InstructionPointer int
	End                bool
}

func NewAmbient() *Ambient {
	return &Ambient{
		Stack:              make([]int, 0),
		Instructions:       make([]common.Instruction, 0),
		InstructionPointer: 0,
		End:                false,
	}
}

func (a *Ambient) LoadProgram(program []common.Instruction) {
	a.Instructions = program
}

func (a *Ambient) Run() common.Error {
	if a.InstructionPointer < 0 || a.InstructionPointer >= len(a.Instructions) {
		return common.IllegalInstruction
	}

	instruction := a.Instructions[a.InstructionPointer]

	switch instruction.Type {
	case common.Push:
		// Push a value onto the stack.
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. PSH 1
		//		2. PRINT_STACK: [0, 1, 1]

		a.Stack = append(a.Stack, instruction.Operand)
		a.InstructionPointer++
		break

	case common.Duplicate:
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
		break

	case common.Plus:
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
		break

	case common.Minus:
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
		break

	case common.Multiply:
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
		break

	case common.Divide:
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
		break

	case common.Jump:
		// Jump to a new instruction.
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. JMP 2
		// 		2. PRINT_STACK: [0, 1]
		if instruction.Operand < 0 || instruction.Operand >= len(a.Instructions) {
			return common.IllegalInstruction
		}

		a.InstructionPointer = instruction.Operand
		break

	case common.JumpIfTrue:
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
		break

	case common.Equal:
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
		break

	case common.End:
		a.End = true
		a.InstructionPointer++
		break

	default:
		return common.IllegalInstruction
	}

	return common.Ok
}

func (a *Ambient) PrintStack() {
	fmt.Println("Stack:")
	if len(a.Stack) == 0 {
		fmt.Println("	Stack is empty")
	}

	for i, v := range a.Stack {
		fmt.Printf("	%d: %d\n", i, v)
	}

	fmt.Println()
}