package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vpaulo/seda/evaluator"
	"github.com/vpaulo/seda/lexer"
	"github.com/vpaulo/seda/object"
	"github.com/vpaulo/seda/parser"
	"github.com/vpaulo/seda/pkg"
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

	// Set command line arguments for OS.args()
	evaluator.SetCommandLineArgs(flag.Args())

	if *help_flag {
		usage()
		os.Exit(0)
	}

	// If no file is provided, start REPL
	if flag.NArg() == 0 {
		start_REPL()
		return
	}

	// Check if first argument is a package management command
	command := flag.Arg(0)
	switch command {
	case "install":
		handle_install()
		return
	case "update":
		handle_update()
		return
	case "remove":
		handle_remove()
		return
	case "list":
		handle_list()
		return
	}

	// Otherwise, treat as a source file
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
		for _, err := range p.FormatErrors() {
			fmt.Fprintf(os.Stderr, "  %s\n", err)
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
		// Set the source directory for module resolution
		abs_path, _ := filepath.Abs(filename)
		env.SourceDir = filepath.Dir(abs_path)
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
	// Set the source directory for module resolution
	abs_path, _ := filepath.Abs(filename)
	env.SourceDir = filepath.Dir(abs_path)
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
	fmt.Println("If no source file is provided, starts an interactive REPL.")
	fmt.Println()
	fmt.Println("OPTIONS:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("PACKAGE MANAGEMENT COMMANDS:")
	fmt.Println("  seda install <package-url>   # Install a package from git repository")
	fmt.Println("  seda update <package-name>   # Update an installed package")
	fmt.Println("  seda remove <package-name>   # Remove an installed package")
	fmt.Println("  seda list                    # List all installed packages")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  seda                                      # Start interactive REPL")
	fmt.Println("  seda program.s                            # Execute program.s")
	fmt.Println("  seda -test program.s                      # Run tests in program.s")
	fmt.Println("  seda -ast program.s                       # Show AST of program.s")
	fmt.Println("  seda -verbose program.s                   # Execute with detailed output")
	fmt.Println("  seda -help                                # Show this help message")
	fmt.Println("  seda install github.com/user/awesome-lib  # Install a package")
	fmt.Println("  seda list                                 # List installed packages")
}

const PROMPT = ">> "

func start_REPL() {
	fmt.Println("Welcome to the Seda Language REPL!")
	fmt.Println("Type 'exit' or 'quit' to exit, 'help' for help")
	fmt.Println()

	env := object.NewEnvironment()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(PROMPT)

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())

		// Handle special commands
		switch line {
		case "exit", "quit":
			fmt.Println("Goodbye!")
			return
		case "help":
			printReplHelp()
			continue
		case "clear":
			clearEnvironment(env)
			fmt.Println("Environment cleared")
			continue
		case "":
			continue
		}

		// Parse and evaluate input
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if p.HasErrors() {
			fmt.Println("Parser errors:")
			for _, err := range p.FormatErrors() {
				fmt.Printf("  %s\n", err)
			}
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			fmt.Printf("%s\n", evaluated.Inspect())
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

func printReplHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  help    - Show this help message")
	fmt.Println("  clear   - Clear all variables and functions")
	fmt.Println("  exit    - Exit the REPL")
	fmt.Println("  quit    - Exit the REPL")
	fmt.Println()
	fmt.Println("Language features:")
	fmt.Println("  Variables: var x = 5")
	fmt.Println("  Functions: fn add(a, b) :: return a + b end")
	fmt.Println("  Arrays: [1, 2, 3]")
	fmt.Println("  Maps: {\"key\": \"value\"}")
	fmt.Println("  Control flow: if, else, case, for")
	fmt.Println("  Testing: check \"test\" :: x is 5 end")
	fmt.Println()
}

func clearEnvironment(env *object.Environment) {
	// Create a new clean environment
	*env = *object.NewEnvironment()
}

// Package management commands

func handle_install() {
	if flag.NArg() < 2 {
		fmt.Println("Error: package URL required")
		fmt.Println("Usage: seda install <package-url>")
		os.Exit(1)
	}
	package_url := flag.Arg(1)
	manager := pkg.NewManager()
	if err := manager.Install(package_url); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func handle_update() {
	if flag.NArg() < 2 {
		fmt.Println("Error: package name required")
		fmt.Println("Usage: seda update <package-name>")
		os.Exit(1)
	}
	package_name := flag.Arg(1)
	manager := pkg.NewManager()
	if err := manager.Update(package_name); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func handle_remove() {
	if flag.NArg() < 2 {
		fmt.Println("Error: package name required")
		fmt.Println("Usage: seda remove <package-name>")
		os.Exit(1)
	}
	package_name := flag.Arg(1)
	manager := pkg.NewManager()
	if err := manager.Remove(package_name); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func handle_list() {
	manager := pkg.NewManager()
	if err := manager.List(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
