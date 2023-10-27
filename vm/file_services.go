package vm

/*
func (a *VirtualMachine) parseNaiveInstructionsFromTokens(asm []token.Token) {
	if len(asm) == 0 {
		color.Set(color.FgHiRed)
		defer color.Unset()

		log.Fatal("No tokens!")
	}

	for i := 0; i < len(asm); i++ {
		switch asm[i].Kind {
		case token.Sum:
			a.Instructions = append(a.Instructions, NewPlus())

		case token.Subtract:
			a.Instructions = append(a.Instructions, NewMinus())

		case token.Multiply:
			a.Instructions = append(a.Instructions, NewMultiply())

		case token.Divide:
			a.Instructions = append(a.Instructions, NewDivide())

		case token.Jump:
			if i+1 >= len(asm) {
				color.Set(color.FgHiRed)
				defer color.Unset()

				// TODO(jejikeh): Maybe print token info
				log.Fatal("Jump instruction is out of bounds.")
			}

			operand := asm[i+1]
			if operand.Kind != token.Identifier && operand.Kind != token.Number {
				color.Set(color.FgHiRed)
				defer color.Unset()

				log.Fatal("Jump operand is not a number or identifier.")
			}

			if operand.IntegerValue < 0 {
				color.Set(color.FgHiRed)
				defer color.Unset()

				log.Fatalf("Jump operand is negative: [%d]! Probably, it was not resolved.\n", operand.IntegerValue)
			}

			a.Instructions = append(a.Instructions, NewJump(operand.IntegerValue))
			i++

		case token.Push:
			if i+1 >= len(asm) {
				color.Set(color.FgHiRed)
				defer color.Unset()

				log.Fatal("Push instruction is out of bounds.")
			}

			operand := asm[i+1]
			if operand.Kind != token.Number {
				color.Set(color.FgHiRed)
				defer color.Unset()

				log.Fatalf("Push operand is not a number: [%s]\n", operand.Kind)
			}

			a.Instructions = append(a.Instructions, NewPush(operand.IntegerValue))
			i++

		case token.Duplicate:
			// TODO: Move it to separate
			if i+1 >= len(asm) {
				color.Set(color.FgHiRed)
				defer color.Unset()

				log.Fatal("Duplicate instruction is out of bounds.")
			}

			operand := asm[i+1]
			if operand.Kind != token.Number {
				color.Set(color.FgHiRed)
				defer color.Unset()

				log.Fatal("Duplicate operand is not a number.")
			}

			a.Instructions = append(a.Instructions, NewDuplicate(operand.IntegerValue))
			i++

		case token.EndOfLine:
			continue

		default:
			color.Set(color.FgHiRed)
			defer color.Unset()

			log.Printf("Unknown instruction: [%s]\n", asm[i].Kind)
		}
	}
}
*/
