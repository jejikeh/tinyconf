using System.Text;
using Ambient.Models;

namespace Ambient.Services;

public class InstructionSaveLoaderService
{
    private readonly VirtualMachine _virtualMachine;
    private readonly InstructionSerializer _instructionSerializer;
    private readonly AssemblerTranslatorService _assemblerTranslatorService;

    public InstructionSaveLoaderService(VirtualMachine virtualMachine, InstructionSerializer instructionSerializer, AssemblerTranslatorService assemblerTranslatorService)
    {
        _virtualMachine = virtualMachine;
        _instructionSerializer = instructionSerializer;
        _assemblerTranslatorService = assemblerTranslatorService;
    }

    public void SaveBinary(string path)
    {
        using var stream = new FileStream(path, FileMode.Append);
        var bytes = _instructionSerializer.SerializeInstructions(_virtualMachine.Instructions.ToArray());
        stream.Write(bytes, 0, bytes.Length);
    }
    
    public void LoadBinary(string path)
    {
        using var stream = new FileStream(path, FileMode.Open);
        var bytes = new byte[stream.Length];
        _ = stream.Read(bytes, 0, bytes.Length);
        _virtualMachine.Instructions = _instructionSerializer.DeserializeInstruction(bytes).ToList();
    }

    public void LoadAssembly(string path)
    {
        var lines = File.ReadLines(path);
        foreach (var line in lines)
        {
            var instruction = _assemblerTranslatorService.TranslateLine(line);
            _virtualMachine.Instructions.Add(instruction);
        }
    }
}