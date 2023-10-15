namespace Ambient.Models;

[Serializable]
public class Instruction
{
    public enum InstructionType
    {
        Push,
        Duplicate,
        Plus,
        Minus,
        Multiply,
        Divide,
        Jump,
        JumpIfTrue,
        Equal,
        End,
    }
    
    public InstructionType Type { get; set; }
    public int Operand { get; set; }

    public static InstructionType GetInstructionType(string instruction)
    {
        return instruction switch
        {
            "psh" => InstructionType.Push,
            "dplc" => InstructionType.Duplicate,
            "sum" => InstructionType.Plus,
            "sub" => InstructionType.Minus,
            "mul" => InstructionType.Multiply,
            "div" => InstructionType.Divide,
            "jmp" => InstructionType.Jump,
            "jif" => InstructionType.JumpIfTrue,
            "eq" => InstructionType.Equal,
            _ => InstructionType.End
        };
    }
}