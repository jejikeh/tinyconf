package common

//go:generate stringer -type=Error
type Error int

const (
	Ok Error = iota
	StackOverflow
	StackUnderflow
	IllegalInstruction
	IllegalInstructionAccess
	DivisionByZero
	UnknownOperand
)

//go:generate stringer -type=InstructionType
type InstructionType int

const (
	Push InstructionType = iota
	Duplicate
	Plus
	Minus
	Multiply
	Divide
	Jump
	JumpIfTrue
	Equal
	End
)

type Instruction struct {
	Type    InstructionType
	Operand int
}

func NewInstruction(type_ InstructionType, operand int) Instruction {
	return Instruction{
		Type:    type_,
		Operand: operand,
	}
}

func NewPush(operand int) Instruction {
	return NewInstruction(Push, operand)
}

func NewDuplicate(operand int) Instruction {
	return NewInstruction(Duplicate, operand)
}

func NewPlus() Instruction {
	return NewInstruction(Plus, 0)
}

func NewMinus() Instruction {
	return NewInstruction(Minus, 0)
}

func NewMultiply() Instruction {
	return NewInstruction(Multiply, 0)
}

func NewDivide() Instruction {
	return NewInstruction(Divide, 0)
}

func NewJump(operand int) Instruction {
	return NewInstruction(Jump, operand)
}

func NewJumpIfTrue(operand int) Instruction {
	return NewInstruction(JumpIfTrue, operand)
}

func NewEqual() Instruction {
	return NewInstruction(Equal, 0)
}

func NewEnd() Instruction {
	return NewInstruction(End, 0)
}

var AmbientAsmInstructionType = map[string]InstructionType{
	"psh":  Push,
	"plus": Plus,
	"dplc": Duplicate,
	"jmp":  Jump,
}
