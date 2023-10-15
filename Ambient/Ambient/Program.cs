using Ambient.Models;
using Ambient.Services;
using CommandLine;
using Error = Ambient.Models.Error;

internal class Program
{
    public class InstallCommand
    {
        public class Options
        {
            [Option('i', "input", Required = true, HelpText = "Input file.")]
            public string Input { get; set; }
            
            [Option('o', "output", Required = true, HelpText = "Output file.")]
            public string Output { get; set; }
        }
    }
    
    public static void Main(string[] args)
    {
        var vm = new VirtualMachine();
        var ats = new AssemblerTranslatorService();
        var isls = new InstructionSaveLoaderService(vm, new InstructionSerializer(), ats);
        var irs = new VirtualMachineInstructionRunner(vm);
        
        Parser.Default.ParseArguments<InstallCommand.Options>(args)
            .WithParsed<InstallCommand.Options>(o =>
            {
                // isls.LoadAssembly(o.Input);
                isls.LoadBinary(o.Input);
                
                while (!vm.End)
                {
                    var error = irs.ExecuteInstructions();
                    if (error == Error.Ok)
                    {
                        continue;
                    }

                    Console.WriteLine(error);
                    break;
                }

                irs.PrintStack();
                isls.SaveBinary(o.Output);
            });
    }
}