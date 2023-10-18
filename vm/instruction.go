package vm

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	"psh":   Push,
	"dplc":  Duplicate,
	"sum":   Plus,
	"minus": Minus,
	"mul":   Multiply,
	"div":   Divide,
	"jmp":   Jump,
	"jif":   JumpIfTrue,
	"eq":    Equal,
	"end":   End,
}

var AmbientAsmInstruction = map[InstructionType]string{
	Push:       "psh",
	Duplicate:  "dplc",
	Plus:       "sum",
	Minus:      "minus",
	Multiply:   "mul",
	Divide:     "div",
	Jump:       "jmp",
	JumpIfTrue: "jif",
	Equal:      "eq",
	End:        "end",
}

func (a *VirtualMachine) DumpDisasembleInstructions(s *string) {
	for _, v := range a.Instructions {
		instructionTextRepresentation := AmbientAsmInstruction[v.Type]
		if v.Type == Plus || v.Type == Minus || v.Type == Multiply || v.Type == Divide || v.Type == End {
			*s += instructionTextRepresentation + "\n"
			continue
		}

		instructionTextRepresentation += " " + strconv.Itoa(v.Operand)
		*s += instructionTextRepresentation + "\n"
	}
}

func (a *VirtualMachine) DumpDisasembleInstructionsToFile(outputPath string) {
	err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(outputPath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for _, v := range a.Instructions {
		instructionTextRepresentation := AmbientAsmInstruction[v.Type]
		if v.Type == Plus || v.Type == Minus || v.Type == Multiply || v.Type == Divide || v.Type == End {
			f.Write([]byte(instructionTextRepresentation))
			f.Write([]byte("\n"))
			continue
		}

		instructionTextRepresentation += " " + strconv.Itoa(v.Operand)
		f.Write([]byte(instructionTextRepresentation))
		f.Write([]byte("\n"))
	}
}

func (a *VirtualMachine) SaveProgramToNewFile(outputPath string) {
	err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(outputPath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.Write(a.serializeInstructions())
	if err != nil {
		log.Fatal(err)
	}
}

func (a *VirtualMachine) LoadProgramFromFile(filepath string) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Err")
	}

	a.LoadProgram(deserializeInstructions(content))
}

// TODO: Maybe change return type to []byte, error?
func (a *VirtualMachine) serializeInstructions() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(a.Instructions)

	if err != nil {
		log.Fatal("Error encoding instructions: ", err)
	}

	return buff.Bytes()
}

func deserializeInstructions(buff []byte) []Instruction {
	var instructions []Instruction

	dec := gob.NewDecoder(bytes.NewBuffer(buff))
	err := dec.Decode(&instructions)
	if err != nil {
		log.Fatal("Error decoding instructions: ", err)
	}

	return instructions
}

func (a *VirtualMachine) translateByteCodeLineToInstruction(line string) (Instruction, error) {
	if len(line) == 0 {
		return NewEnd(), nil
	}

	if line[0] == '#' {
		return NewEnd(), errors.New("Comment")
	}

	if line[0] == ':' {
		a.Labels[line[1:]] = len(a.Instructions)
		return NewEnd(), errors.New("Label")
	}

	instructionByDelimiter := strings.Split(line, " ")

	instKind, ok := AmbientAsmInstructionType[instructionByDelimiter[0]]
	if !ok {
		log.Fatalf("Unknown instruction: [%s]\n", instructionByDelimiter[0])
	}

	if len(instructionByDelimiter) < 2 {
		return NewInstruction(instKind, 0), nil
	}

	operand, err := strconv.Atoi(instructionByDelimiter[1])

	if (instKind == Jump || instKind == JumpIfTrue) && err != nil {
		if val, ok := a.Labels[instructionByDelimiter[1]]; ok {
			return NewInstruction(instKind, val), nil
		}

		log.Printf("Unknown label: [%s]\n", instructionByDelimiter[1])
		a.NotResolvedLabels[instructionByDelimiter[1]] = len(a.Instructions)
	}

	if err != nil {
		log.Fatalf("Invalid operand: %s\n", instructionByDelimiter[1])
	}

	return NewInstruction(instKind, operand), nil
}
