using Ambient.Models;

namespace Ambient.Services;

public class AssemblerTranslatorService
{
    public Instruction TranslateLine(string line)
    {
        if (!line.Any())
        {
            return new Instruction()
            {
                Type = Instruction.InstructionType.End
            };
            
        }
        
        var instructionDel = line.Split(" ");
        var instructionKind = Instruction.GetInstructionType(instructionDel[0]);

        if (instructionDel.Length != 2)
        {
            return new Instruction
            {
                Type = instructionKind,
                Operand = 0
            };
        }
        
        return new Instruction
        {
            Type = instructionKind,
            Operand = int.Parse(instructionDel[1])
        };
    }
}