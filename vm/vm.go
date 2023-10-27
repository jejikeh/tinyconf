package vm

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/jejikeh/ambient/common"
	"github.com/jejikeh/ambient/lexer"
	"github.com/jejikeh/ambient/token"
)

type VirtualMachine struct {
	Stack              []int
	Instructions       []token.Token
	Labels             map[string]int
	NotResolvedLabels  map[string]int
	InstructionPointer int
}

func NewVirtualMachine() *VirtualMachine {
	return &VirtualMachine{
		Stack:              make([]int, 0),
		Instructions:       make([]token.Token, 0),
		Labels:             make(map[string]int),
		NotResolvedLabels:  make(map[string]int),
		InstructionPointer: 0,
	}
}

func (a *VirtualMachine) LoadNaiveFromSourceFile(sourcePath string) {
	// TODO(jejikeh): fix this allocation
	l := lexer.NewLexerFromSource(sourcePath)
	a.LoadProgram(l.Tokenize())
}

func (a *VirtualMachine) LoadNaiveFromSourceBinary(sourcePath string) {
	// TODO(jejikeh): fix this allocation
	l := lexer.NewLexerFromBinary(sourcePath)
	a.LoadProgram(l.Tokens)
}

func (a *VirtualMachine) LoadProgram(program []token.Token) {
	a.Instructions = program
}

func (a *VirtualMachine) Run() common.Error {
	if a.InstructionPointer < 0 || a.InstructionPointer >= len(a.Instructions) {
		return common.IllegalInstruction
	}

	instruction := a.Instructions[a.InstructionPointer]

	switch instruction.Kind {
	case token.Push:
		// Push a value onto the stack.
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. PSH 1
		//		2. PRINT_STACK: [0, 1, 1]

		a.Stack = append(a.Stack, a.Instructions[a.InstructionPointer+1].IntegerValue)
		a.InstructionPointer++

	case token.Duplicate:
		// Duplicate the top of the stack.
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. DPLC 0
		//		2. PRINT_STACK: [0, 1, 0]

		if len(a.Stack)-a.Instructions[a.InstructionPointer+1].IntegerValue <= 0 {
			return common.StackUnderflow
		}

		if instruction.IntegerValue < 0 {
			return common.IllegalInstruction
		}

		a.Stack = append(a.Stack, a.Stack[len(a.Stack)-1-a.Instructions[a.InstructionPointer+1].IntegerValue])
		a.InstructionPointer++

	case token.Sum:
		// Add the top two values on the stack.
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. PSH 1
		// 		2. PSH 1
		// 		3. SUM
		//		4. PRINT_STACK: [0, 1, 2]

		if len(a.Stack) < 2 {
			return common.StackUnderflow
		}

		a.Stack[len(a.Stack)-2] = a.Stack[len(a.Stack)-2] + a.Stack[len(a.Stack)-1]
		a.Stack = a.Stack[:len(a.Stack)-1]
		a.InstructionPointer++

	case token.Subtract:
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

	case token.Multiply:
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

	case token.Divide:
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

	case token.Jump:
		// Jump to a new instruction.
		// EXAMPLE:
		// 		0. PRINT_STACK: [0, 1]
		// 		1. JMP 2
		// 		2. PRINT_STACK: [0, 1]
		if a.Instructions[a.InstructionPointer+1].IntegerValue < 0 || a.Instructions[a.InstructionPointer+1].IntegerValue >= len(a.Instructions) {
			return common.IllegalInstructionAccess
		}

		a.InstructionPointer = a.Instructions[a.InstructionPointer+1].IntegerValue

	case token.JumpIfTrue:
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

		a.InstructionPointer = a.Instructions[a.InstructionPointer+1].IntegerValue

	case token.Equal:
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

	default:
		log.Printf("Unknown instruction: [%s]\n", instruction.Kind)
		a.InstructionPointer++
	}

	return common.Ok
}

func (a *VirtualMachine) Execute(executingLimit int, printCurrentInstruction bool) {
	isInfinite := executingLimit < 0
	for i := 0; (i < executingLimit && (a.Instructions[a.InstructionPointer].Kind != token.EndOfLine)) || isInfinite; i++ {
		err := a.Run()
		if err != common.Ok {
			color.Set(color.FgHiRed)
			defer color.Unset()

			log.Printf("Error: %s\n", err.String())
			a.PrintStack()
			panic(1)
		}

		if printCurrentInstruction {
			log.Printf("[%d] Current pointer -> [%d: %s]\n", i, a.InstructionPointer, a.Instructions[a.InstructionPointer].Kind)
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
		fmt.Printf("	%d: %s %d\n", i, v.Kind, v.IntegerValue)
	}

	fmt.Println()
}
