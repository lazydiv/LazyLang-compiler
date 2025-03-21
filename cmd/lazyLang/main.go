package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/lazydiv/lazyLang-compiler/internal/codegen"
	"github.com/lazydiv/lazyLang-compiler/internal/lexer"
	"github.com/lazydiv/lazyLang-compiler/internal/parser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: lazylang <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	l := lexer.NewLexer(string(source))
	p := parser.NewParser(l)
	program := p.ParseProgram()

	cg := codegen.NewCodeGen()
	goCode := cg.Generate(program)

	outFile := strings.TrimSuffix(filename, ".lazy") + ".go"
	err = os.WriteFile(outFile, []byte(goCode), 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Compiled %s to %s\n", filename, outFile)
	// run the output file
	cmd := exec.Command("go", "run", outFile)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running output file: %v\n", err)
		os.Exit(1)
	}

}
