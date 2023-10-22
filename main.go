package main

import (
	"flag"
	"fmt"

	"github.com/jejikeh/ambient/lexer"
	"github.com/jejikeh/ambient/vm"
)

func main() {
	debugFlag := flag.Bool("debug", false, "Debug")

	sourcePath := flag.String("i", "", "Source file")
	outputPath := flag.String("o", "", "Output file")

	// Build Command
	buildCommand := flag.Bool("build", false, "Build binary")
	defer buildBinary(buildCommand, sourcePath, outputPath, debugFlag)

	// Disassemble Command
	disassembleCommand := flag.Bool("dis", false, "Disassemble binary")
	defer dissembleBinary(disassembleCommand, sourcePath, outputPath)

	// Run Command
	runCommand := flag.Bool("run", false, "Run binary")
	binaryFlag := flag.Bool("x", false, "Binary flag")
	defer runBinary(runCommand, binaryFlag, sourcePath, debugFlag)

	// Lexer Command
	lexerFlag := flag.Bool("lex", false, "Lexer file")
	defer lexerFile(lexerFlag, sourcePath)

	flag.Parse()
}

func dissembleBinary(disassembleFlag *bool, source *string, output *string) {
	if !*disassembleFlag {
		return
	}

	ambient := vm.NewVirtualMachine()
	ambient.LoadProgramFromFile(*source)

	if *output == "" {
		dissasembleContent := ""
		ambient.DumpDisasembleInstructions(&dissasembleContent)
		fmt.Println(dissasembleContent)
		return
	}

	ambient.DumpDisasembleInstructionsToFile(*output)
}

func runBinary(runFlag *bool, binaryFlag *bool, source *string, debug *bool) {
	if !*runFlag {
		return
	}

	ambient := vm.NewVirtualMachine()

	if *binaryFlag {
		ambient.LoadProgramFromFile(*source)
	} else {
		ambient.LoadByteCodeAsmFromFile(*source)
	}

	if *debug {
		ambient.PrintInstructions()
		ambient.Execute(100, true)
		ambient.PrintStack()
		return
	}

	ambient.Execute(100, false)
}

func buildBinary(binaryFlag *bool, source *string, output *string, debug *bool) {
	if !*binaryFlag {
		return
	}

	ambient := vm.NewVirtualMachine()
	ambient.LoadByteCodeAsmFromFile(*source)

	if *debug {
		ambient.PrintInstructions()
	}

	ambient.SaveProgramToNewFile(*output)
}

func lexerFile(lexerFlag *bool, source *string) {
	if !*lexerFlag {
		return
	}

	l := lexer.NewLexerFromFile(*source)
	tokens := l.Tokenize()
	lexer.PrintDebugTokens(tokens)
}
