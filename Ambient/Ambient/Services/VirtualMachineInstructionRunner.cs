using Ambient.Models;

namespace Ambient.Services;

public class VirtualMachineInstructionRunner
{
    private readonly VirtualMachine _virtualMachine;

    public VirtualMachineInstructionRunner(VirtualMachine virtualMachine)
    {
        _virtualMachine = virtualMachine;
    }

    public Error ExecuteInstructions()
    {
	    if (_virtualMachine.Pointer < 0 || _virtualMachine.Pointer >= _virtualMachine.Instructions.Count)
	    {
		    return Error.IllegalInstruction;
	    }

	    var instruction = _virtualMachine.Instructions[_virtualMachine.Pointer];

	    switch (instruction.Type)
	    {
		    case Instruction.InstructionType.Push:
			    // Push a value onto the stack.
			    // EXAMPLE:
			    // 		0. PRINT_STACK: [0, 1]
			    // 		1. PSH 1
			    //		2. PRINT_STACK: [0, 1, 1]

			    _virtualMachine.Stack.Add(instruction.Operand);
			    _virtualMachine.Pointer++;
			    break;

		    case Instruction.InstructionType.Duplicate:
			    // Duplicate the top of the stack.
			    // EXAMPLE:
			    // 		0. PRINT_STACK: [0, 1]
			    // 		1. DPLC 0
			    //		2. PRINT_STACK: [0, 1, 0]

			    if (_virtualMachine.Stack.Count - instruction.Operand < 0)
			    {
				    return Error.StackUnderflow;
			    }

			    if (instruction.Operand < 0)
			    {
				    return Error.IllegalInstruction;
			    }

			    _virtualMachine.Stack.Add(_virtualMachine.Stack[_virtualMachine.Stack.Count - 1 - instruction.Operand]);
			    _virtualMachine.Pointer++;
			    break;

		    case Instruction.InstructionType.Plus:
			    // Add the top two values on the stack.
			    // EXAMPLE:
			    // 		0. PRINT_STACK: [0, 1]
			    // 		1. PSH 1
			    // 		2. PSH 1
			    // 		3. ADD
			    //		4. PRINT_STACK: [0, 1, 2]

			    if (_virtualMachine.Stack.Count < 2)
			    {
				    return Error.StackUnderflow;
			    }

			    _virtualMachine.Stack[^2] += _virtualMachine.Stack[^1];
			    _virtualMachine.Stack = _virtualMachine.Stack.GetRange(0, _virtualMachine.Stack.Count - 1);
			    _virtualMachine.Pointer++;

			    break;

		    case Instruction.InstructionType.Minus:
			    // Subtract the top two values on the stack.
			    // EXAMPLE:
			    // 		0. PRINT_STACK: [0, 1]
			    // 		1. PSH 2
			    // 		2. PSH 1
			    // 		3. SUB
			    // 		4. PRINT_STACK: [0, 1, 1]

			    if (_virtualMachine.Stack.Count < 2)
			    {
				    return Error.StackUnderflow;
			    }

			    _virtualMachine.Stack[^2] -= _virtualMachine.Stack[^1];
			    _virtualMachine.Stack = _virtualMachine.Stack.GetRange(0, _virtualMachine.Stack.Count - 1);
			    _virtualMachine.Pointer++;
			    break;
		    
		    case Instruction.InstructionType.Multiply:
			    // Multiply the top two values on the stack.
			    // EXAMPLE:
			    // 		0. PRINT_STACK: [0, 1]
			    // 		1. PSH 2
			    // 		2. PSH 2
			    // 		3. MUL
			    // 		4. PRINT_STACK: [0, 1, 4]

			    if (_virtualMachine.Stack.Count < 2)
			    {
				    return Error.StackUnderflow;
			    }

			    _virtualMachine.Stack[^2] *= _virtualMachine.Stack[^1];
			    _virtualMachine.Stack = _virtualMachine.Stack.GetRange(0, _virtualMachine.Stack.Count - 1);
			    _virtualMachine.Pointer++;
			    break;

		    case Instruction.InstructionType.Divide:
			    // Divide the top two values on the stack.
			    // EXAMPLE:
			    // 		0. PRINT_STACK: [0, 1]
			    // 		1. PSH 4
			    // 		2. PSH 2
			    // 		3. DIV
			    // 		4. PRINT_STACK: [0, 1, 2]

			    if (_virtualMachine.Stack.Count < 2)
			    {
				    return Error.StackUnderflow;
			    }

			    _virtualMachine.Stack[^2] /= _virtualMachine.Stack[^1];
			    _virtualMachine.Stack = _virtualMachine.Stack.GetRange(0, _virtualMachine.Stack.Count - 1);
			    _virtualMachine.Pointer++;
			    break;

		    case Instruction.InstructionType.Jump:
			    // Jump to a new instruction.
			    // EXAMPLE:
			    // 		0. PRINT_STACK: [0, 1]
			    // 		1. JMP 2
			    // 		2. PRINT_STACK: [0, 1]
			    if (instruction.Operand < 0 || instruction.Operand >= _virtualMachine.Instructions.Count)
			    {
				    return Error.IllegalInstruction;
			    }

			    _virtualMachine.Pointer = instruction.Operand;
			    break;

		    case Instruction.InstructionType.JumpIfTrue:
			    // Jump to a new instruction if the top of the stack is true (1).
			    // EXAMPLE:
			    // 		0. PRINT_STACK: [0, 1]
			    // 		1. PSH 0
			    // 		2. PSH 1
			    // 		3. JMPIF 4
			    // 		4. PRINT_STACK: [0, 1]

			    if (_virtualMachine.Stack.Count < 1)
			    {
				    return Error.StackUnderflow;
			    }

			    if (_virtualMachine.Stack[^1] != 1)
			    {
				    _virtualMachine.Pointer++;
				    break;
			    }

			    _virtualMachine.Stack = _virtualMachine.Stack.GetRange(0, _virtualMachine.Stack.Count - 1);
			    _virtualMachine.Pointer = instruction.Operand;
			    break;

		    case Instruction.InstructionType.Equal:
			    // instruction if the top of the stack is equal.
			    // EXAMPLE:
			    // 		0. PRINT_STACK: [0, 1]
			    // 		1. PSH 1
			    // 		2. PSH 1
			    // 		3. EQ
			    // 		4. PRINT_STACK: [0, 1, 1]
			    if (_virtualMachine.Stack.Count < 2)
			    {
				    return Error.StackUnderflow;
			    }

			    _virtualMachine.Stack[^2] = _virtualMachine.Stack[^2] == _virtualMachine.Stack[^1] ? 1 : 0;
			    _virtualMachine.Stack = _virtualMachine.Stack.GetRange(0, _virtualMachine.Stack.Count - 1);
			    _virtualMachine.Pointer++;
			    break;

		    case Instruction.InstructionType.End:
			    _virtualMachine.End = true;
			    _virtualMachine.Pointer++;
			    break;

		    default:
			    return Error.IllegalInstruction;
	    }
	    
	    return Error.Ok;
    }

    public void PrintStack()
    {
	    Console.WriteLine("Stack: ");

	    if (_virtualMachine.Stack.Count == 0)
	    {
		    Console.WriteLine("\tEmpty");
	    }

	    for (var i = 0; i < _virtualMachine.Stack.Count; i++)
	    {
		    Console.WriteLine($"\t{i}: {_virtualMachine.Stack[i]}");
	    }
	    
	    Console.WriteLine();
    }
    
    public void PrintInstructions()
    {
	    Console.WriteLine("Instructions: ");

	    if (_virtualMachine.Instructions.Count == 0)
	    {
		    Console.WriteLine("\tEmpty");
	    }

	    for (var i = _virtualMachine.Instructions.Count - 1; i >= 0; i--)
	    {
		    Console.WriteLine($"\t{i}: {_virtualMachine.Instructions[i].Type} {_virtualMachine.Instructions[i].Operand}");
	    }
	    
	    Console.WriteLine();
    }
    
    
}