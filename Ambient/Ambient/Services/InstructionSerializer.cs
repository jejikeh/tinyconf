using System.Runtime.Serialization.Formatters.Binary;
using Ambient.Models;

namespace Ambient.Services;

public class InstructionSerializer
{
    private readonly BinaryFormatter _formatter = new BinaryFormatter();

    [Obsolete("Obsolete")]
    public byte[] SerializeInstruction(Instruction instruction)
    {
        using var stream = new MemoryStream();
        _formatter.Serialize(stream, instruction);
        return stream.ToArray();
    }
    
    [Obsolete("Obsolete")]
    public byte[] SerializeInstructions(Instruction[] instruction)
    {
        using var stream = new MemoryStream();
        _formatter.Serialize(stream, instruction);
        return stream.ToArray();
    }

    [Obsolete("Obsolete")]
    public Instruction[] DeserializeInstruction(byte[] data)
    {
        using var memStream = new MemoryStream();
        memStream.Write(data, 0, data.Length);
        memStream.Seek(0, SeekOrigin.Begin);
        var obj = _formatter.Deserialize(memStream);
        return (Instruction[])obj;
    }
}