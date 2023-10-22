package vm

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

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
		color.Set(color.FgHiRed)
		defer color.Unset()

		log.Fatalf("Error reading file: [%v]\n", err)
	}

	a.LoadProgram(deserializeInstructions(content))
}

// TODO: Maybe change return type to []byte, error?
func (a *VirtualMachine) serializeInstructions() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(a.Instructions)

	if err != nil {
		color.Set(color.FgHiRed)
		defer color.Unset()

		log.Fatal("Error encoding instructions: ", err)
	}

	return buff.Bytes()
}

func deserializeInstructions(buff []byte) []Instruction {
	var instructions []Instruction

	dec := gob.NewDecoder(bytes.NewBuffer(buff))
	err := dec.Decode(&instructions)
	if err != nil {
		color.Set(color.FgHiRed)
		defer color.Unset()

		log.Fatal("Error decoding instructions: ", err)
	}

	return instructions
}

var InstructionToString = map[InstructionType]string{
	Push:       "psh",
	Duplicate:  "dupl",
	Plus:       "pls",
	Minus:      "sub",
	Multiply:   "mul",
	Divide:     "div",
	Jump:       "jmp",
	JumpIfTrue: "jift",
	Equal:      "eq",
	End:        "end",
}
