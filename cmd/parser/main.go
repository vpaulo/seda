package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vpaulo/seda/evaluator"
	"github.com/vpaulo/seda/lexer"
	"github.com/vpaulo/seda/object"
	"github.com/vpaulo/seda/parser"
)

var (
	test_mode    = flag.Bool("test", false, "Run tests instead of executing code")
	ast_mode     = flag.Bool("ast", false, "Show AST and exit (don't execute)")
	verbose_mode = flag.Bool("verbose", false, "Show detailed execution information")
	help_flag    = flag.Bool("help", false, "Show help message")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if *help_flag {
		usage()
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Error: Please provide exactly one source file\n\n")
		usage()
		os.Exit(1)
	}

	filename := flag.Arg(0)

	// Check file extension
	if ext := filepath.Ext(filename); ext != ".s" {
		fmt.Fprintf(os.Stderr, "Warning: File '%s' doesn't have .s extension\n", filename)
	}

	input, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if *verbose_mode {
		fmt.Printf("File: %s (%d bytes)\n", filename, len(input))
	}

	l := lexer.New(string(input))
	p := parser.New(l)
	program := p.ParseProgram()

	if p.HasErrors() {
		fmt.Fprintf(os.Stderr, "Parse errors in %s:\n", filename)
		for i, err := range p.FormatErrors() {
			fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, err)
		}
		os.Exit(1)
	}

	if *verbose_mode {
		fmt.Printf("Parsed successfully (%d statements)\n", len(program.Statements))
	}

	// AST mode - just show AST and exit
	if *ast_mode {
		fmt.Println("Abstract Syntax Tree:")
		fmt.Println(program.String())
		return
	}

	// Test mode - run tests
	if *test_mode {
		if *verbose_mode {
			fmt.Printf("Running tests in %s...\n", filename)
		} else {
			fmt.Println("Running tests...")
		}
		env := object.NewEnvironment()
		test_result := evaluator.RunTests(program, env)
		fmt.Println(test_result.String())

		// Exit with error code if tests failed
		if test_result.Failed > 0 {
			os.Exit(1)
		}
		return
	}

	// Normal execution mode
	if *verbose_mode {
		fmt.Printf("Executing %s...\n", filename)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil {
		switch result := result.(type) {
		case *object.Error:
			fmt.Fprintf(os.Stderr, "Runtime error: %s\n", result.Message)
			os.Exit(1)
		default:
			if *verbose_mode {
				fmt.Printf("Program completed. Final result: %s\n", result.Inspect())
			}
			// Only print result if it's not null and not a function definition
			if result.Type() != object.NULL_OBJ && result.Type() != object.FUNCTION_OBJ {
				fmt.Println(result.String())
			}
		}
	}
}

func usage() {
	fmt.Printf("Usage: seda [OPTIONS] <source-file>\n\n")
	fmt.Println("A Programming Language Interpreter")
	fmt.Println()
	fmt.Println("OPTIONS:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  seda program.s              # Execute program.s")
	fmt.Println("  seda -test program.s        # Run tests in program.s")
	fmt.Println("  seda -ast program.s         # Show AST of program.s")
	fmt.Println("  seda -verbose program.s     # Execute with detailed output")
}
