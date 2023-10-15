namespace Ambient.Models;

public class VirtualMachine
{
    public List<int> Stack { get; set; } = new List<int>();
    public List<Instruction> Instructions { get; set; } = new List<Instruction>();
    public int Pointer { get; set; }
    public bool End { get; set; }
}