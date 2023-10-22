package vm

/*
func (a *VirtualMachine) LoadByteCodeAsmFromFile(sourcePath string) {
	readFile, err := os.Open(sourcePath)

	if err != nil {
		color.Set(color.FgHiRed)
		defer color.Unset()

		log.Fatalf("error opening file: %v\n", err)
	}

	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		// TODO(jejikeh): fix that allocation
		l := lexer.NewLexer(fileScanner.Text())
		tokens := l.Tokenize()

		a.loadByteCodeAsmFromString(tokens)
	}
}

func (a *VirtualMachine) loadByteCodeAsmFromString(asm []lexer.Token) {
	if len(asm) == 0 {
		color.Set(color.FgHiRed)
		defer color.Unset()

		log.Fatal("ByteCode Asm is empty")
	}

	for i := 0; i < len(asm); i++ {
		switch asm[i].Kind {
		case lexer.Sum:
			a.Instructions = append(a.Instructions, NewPlus())

		case lexer.Sub:
			a.Instructions = append(a.Instructions, NewMinus())

		case lexer.Mul:
			a.Instructions = append(a.Instructions, NewMultiply())

		case lexer.Div:
			a.Instructions = append(a.Instructions, NewDivide())

		case lexer.Jump:
			if i+1 >= len(asm) {
				color.Set(color.FgHiRed)
				defer color.Unset()

				// TODO(jejikeh): Maybe print token info
				log.Fatal("Jump instruction is out of bounds.")
			}

			operand := asm[i+1]
			if operand.Kind != lexer.OperandNumber {
				color.Set(color.FgHiRed)
				defer color.Unset()

				log.Fatal("Jump operand is not a number.")
			}

			numOperand, err := strconv.Atoi(operand.Value)
			if err != nil {
				color.Set(color.FgHiRed)
				defer color.Unset()

				log.Fatal("Jump operand is not a number.")
			}

			a.Instructions = append(a.Instructions, NewJump(numOperand))
			i++

		case lexer.Push:
			if i+1 >= len(asm) {
				color.Set(color.FgHiRed)
				defer color.Unset()

				log.Fatal("Push instruction is out of bounds.")
			}

			operand := asm[i+1]
			if operand.Kind != lexer.OperandNumber {
				color.Set(color.FgHiRed)
				defer color.Unset()

				log.Fatal("Push operand is not a number.")
			}

			numOperand, err := strconv.Atoi(operand.Value)
			if err != nil {
				color.Set(color.FgHiRed)
				defer color.Unset()

				log.Fatal("Push operand is not a number.")
			}

			a.Instructions = append(a.Instructions, NewPush(numOperand))
			i++

		case lexer.Duplicate:
			// TODO: Move it to separate
			if i+1 >= len(asm) {
				color.Set(color.FgHiRed)
				defer color.Unset()

				log.Fatal("Duplicate instruction is out of bounds.")
			}

			operand := asm[i+1]
			if operand.Kind != lexer.OperandNumber {
				color.Set(color.FgHiRed)
				defer color.Unset()

				log.Fatal("Duplicate operand is not a number.")
			}

			numOperand, err := strconv.Atoi(operand.Value)
			if err != nil {
				color.Set(color.FgHiRed)
				defer color.Unset()

				log.Fatal("Duplicate operand is not a number.")
			}

			a.Instructions = append(a.Instructions, NewDuplicate(numOperand))
			i++

		default:
			color.Set(color.FgHiRed)
			defer color.Unset()

			log.Fatalf("Unknown instruction: [%s]\n", asm[i].Value)
		}
	}
}
*/
