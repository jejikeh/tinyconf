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
