package vm

type Error string

const (
	Ok                       = "Ok"
	StackOverflow            = "Stack overflow"
	StackUnderflow           = "Stack underflow"
	IllegalInstruction       = "Illegal instruction"
	IllegalInstructionAccess = "Access to illegal instruction"
	DivisionByZero           = "Division by zero"
	UnknownOperand           = "Unknown operand"
)
