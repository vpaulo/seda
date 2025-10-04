package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vpaulo/seda/lexer"
)

var (
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

	fmt.Printf("Lexer: %v\n", l)
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
