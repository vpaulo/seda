package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// func TestE2EREPL(t *testing.T) {
// 	// Test basic REPL functionality
// 	cmd := exec.Command("go", "run", "cmd/repl/main.go")
// 	cmd.Stdin = strings.NewReader("var x = 5\nx\nexit\n")

// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		t.Fatalf("REPL command failed: %v\nOutput: %s", err, output)
// 	}

// 	output_str := string(output)
// 	if !strings.Contains(output_str, "Welcome to the Programming Language REPL") {
// 		t.Errorf("Expected REPL welcome message, got %q", output_str)
// 	}

// 	if !strings.Contains(output_str, "5") {
// 		t.Errorf("Expected variable value output, got %q", output_str)
// 	}
// }

func TestE2EExamplePrograms(t *testing.T) {
	// Test all example programs
	example_dir := "examples/"

	examples, err := filepath.Glob(filepath.Join(example_dir, "*.s"))
	if err != nil {
		t.Fatalf("Failed to find example files: %v", err)
	}

	if len(examples) == 0 {
		t.Skip("No example files found")
	}

	for _, example := range examples {
		t.Run(filepath.Base(example), func(t *testing.T) {
			// Test normal execution
			cmd := exec.Command("go", "run", "cmd/parser/main.go", "-test", example)

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("Example %s failed to execute: %v\nOutput: %s", example, err, output)
			}

			// Test AST mode
			cmd = exec.Command("go", "run", "cmd/parser/main.go", "-ast", example)
			_, err = cmd.CombinedOutput()
			if err != nil {
				t.Errorf("Example %s failed in AST mode: %v", example, err)
			}

			// Test verbose mode
			cmd = exec.Command("go", "run", "cmd/parser/main.go", "-verbose", example)
			_, err = cmd.CombinedOutput()
			if err != nil {
				t.Errorf("Example %s failed in verbose mode: %v", example, err)
			}
		})
	}
}

func TestE2EHelpAndUsage(t *testing.T) {
	// Test help flag
	cmd := exec.Command("go", "run", "cmd/parser/main.go", "-help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Help command failed: %v", err)
	}

	output_str := string(output)
	if !strings.Contains(output_str, "Usage:") {
		t.Errorf("Expected usage information, got %q", output_str)
	}

	if !strings.Contains(output_str, "OPTIONS:") {
		t.Errorf("Expected options information, got %q", output_str)
	}
}
