package vm

/*
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
		instructionTextRepresentation := InstructionToString[v.Type]
		if v.Kind == token. || v.Type == Minus || v.Type == Multiply || v.Type == Divide || v.Type == End {
			f.Write([]byte(instructionTextRepresentation))
			f.Write([]byte("\n"))
			continue
		}

		instructionTextRepresentation += " " + strconv.Itoa(v.Operand)
		f.Write([]byte(instructionTextRepresentation))
		f.Write([]byte("\n"))
	}
}

func (a *VirtualMachine) DumpDisasembleInstructions(s *string) {
	for _, v := range a.Instructions {
		instructionTextRepresentation := InstructionToString[v.Type]
		if v.Type == Plus || v.Type == Minus || v.Type == Multiply || v.Type == Divide || v.Type == End {
			*s += instructionTextRepresentation + "\n"
			continue
		}

		instructionTextRepresentation += " " + strconv.Itoa(v.Operand)
		*s += instructionTextRepresentation + "\n"
	}
}
*/
