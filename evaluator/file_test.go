package evaluator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vpaulo/seda/object"
)

// File Module Tests

func TestFileReadWrite(t *testing.T) {
	// Create a temporary file
	tmpFile := filepath.Join(os.TempDir(), "seda_test_read_write.txt")
	defer os.Remove(tmpFile)

	// Write content
	writeInput := `File.write("` + tmpFile + `", "Hello, World!")`
	writeResult := testEval(writeInput)

	if writeResult != object.NULL {
		t.Fatalf("write should return NULL. got=%T (%+v)", writeResult, writeResult)
	}

	// Read content back
	readInput := `
		var content, err = File.read("` + tmpFile + `")
		content
	`
	readResult := testEval(readInput)

	str, ok := readResult.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", readResult, readResult)
	}

	if str.Value != "Hello, World!" {
		t.Errorf("wrong content. expected=%s, got=%s", "Hello, World!", str.Value)
	}
}

func TestFileReadLines(t *testing.T) {
	// Create a temporary file with multiple lines
	tmpFile := filepath.Join(os.TempDir(), "seda_test_read_lines.txt")
	defer os.Remove(tmpFile)

	os.WriteFile(tmpFile, []byte("line1\nline2\nline3"), 0644)

	input := `
		var lines, err = File.read_lines("` + tmpFile + `")
		lines
	`
	result := testEval(input)

	arr, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("result is not Array. got=%T (%+v)", result, result)
	}

	expected := []string{"line1", "line2", "line3"}
	if len(arr.Elements) != len(expected) {
		t.Fatalf("array has wrong length. expected=%d, got=%d", len(expected), len(arr.Elements))
	}

	for i, exp := range expected {
		str, ok := arr.Elements[i].(*object.String)
		if !ok {
			t.Fatalf("element %d is not String. got=%T", i, arr.Elements[i])
		}
		if str.Value != exp {
			t.Errorf("element %d: expected=%s, got=%s", i, exp, str.Value)
		}
	}
}

func TestFileAppend(t *testing.T) {
	// Create a temporary file
	tmpFile := filepath.Join(os.TempDir(), "seda_test_append.txt")
	defer os.Remove(tmpFile)

	// Write initial content
	os.WriteFile(tmpFile, []byte("Hello"), 0644)

	// Append more content
	appendInput := `File.append("` + tmpFile + `", ", World!")`
	appendResult := testEval(appendInput)

	if appendResult != object.NULL {
		t.Fatalf("append should return NULL. got=%T (%+v)", appendResult, appendResult)
	}

	// Read back
	content, _ := os.ReadFile(tmpFile)
	if string(content) != "Hello, World!" {
		t.Errorf("wrong content. expected=%s, got=%s", "Hello, World!", string(content))
	}
}

func TestFileDelete(t *testing.T) {
	// Create a temporary file
	tmpFile := filepath.Join(os.TempDir(), "seda_test_delete.txt")
	os.WriteFile(tmpFile, []byte("test"), 0644)

	// Delete it
	deleteInput := `File.delete("` + tmpFile + `")`
	deleteResult := testEval(deleteInput)

	if deleteResult != object.NULL {
		t.Fatalf("delete should return NULL. got=%T (%+v)", deleteResult, deleteResult)
	}

	// Verify it's gone
	if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
		t.Error("file should have been deleted")
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary file
	tmpFile := filepath.Join(os.TempDir(), "seda_test_exists.txt")
	os.WriteFile(tmpFile, []byte("test"), 0644)
	defer os.Remove(tmpFile)

	tests := []struct {
		input    string
		expected bool
	}{
		{`File.exists("` + tmpFile + `")`, true},
		{`File.exists("/nonexistent/file/path.txt")`, false},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testBooleanObject(t, result, tt.expected)
	}
}

func TestFileSize(t *testing.T) {
	// Create a temporary file with known size
	tmpFile := filepath.Join(os.TempDir(), "seda_test_size.txt")
	content := "Hello, World!"
	os.WriteFile(tmpFile, []byte(content), 0644)
	defer os.Remove(tmpFile)

	input := `
		var size, err = File.size("` + tmpFile + `")
		size
	`
	result := testEval(input)

	testNumberObject(t, result, float64(len(content)))
}

func TestFileIsFile(t *testing.T) {
	// Create a temporary file
	tmpFile := filepath.Join(os.TempDir(), "seda_test_is_file.txt")
	os.WriteFile(tmpFile, []byte("test"), 0644)
	defer os.Remove(tmpFile)

	// Create a temporary directory
	tmpDir := filepath.Join(os.TempDir(), "seda_test_is_file_dir")
	os.Mkdir(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		input    string
		expected bool
	}{
		{`File.is_file("` + tmpFile + `")`, true},
		{`File.is_file("` + tmpDir + `")`, false},
		{`File.is_file("/nonexistent/path")`, false},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testBooleanObject(t, result, tt.expected)
	}
}

func TestFileIsDir(t *testing.T) {
	// Create a temporary file
	tmpFile := filepath.Join(os.TempDir(), "seda_test_is_dir_file.txt")
	os.WriteFile(tmpFile, []byte("test"), 0644)
	defer os.Remove(tmpFile)

	// Create a temporary directory
	tmpDir := filepath.Join(os.TempDir(), "seda_test_is_dir")
	os.Mkdir(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		input    string
		expected bool
	}{
		{`File.is_dir("` + tmpDir + `")`, true},
		{`File.is_dir("` + tmpFile + `")`, false},
		{`File.is_dir("/nonexistent/path")`, false},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testBooleanObject(t, result, tt.expected)
	}
}

func TestFileListDir(t *testing.T) {
	// Create a temporary directory with some files
	tmpDir := filepath.Join(os.TempDir(), "seda_test_list_dir")
	os.Mkdir(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	// Create some files
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file3.txt"), []byte("test"), 0644)

	input := `
		var files, err = File.list_dir("` + tmpDir + `")
		files.length()
	`
	result := testEval(input)

	testNumberObject(t, result, 3)
}

func TestFileMkdir(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "seda_test_mkdir")
	defer os.RemoveAll(tmpDir)

	input := `File.mkdir("` + tmpDir + `")`
	result := testEval(input)

	if result != object.NULL {
		t.Fatalf("mkdir should return NULL. got=%T (%+v)", result, result)
	}

	// Verify directory was created
	if info, err := os.Stat(tmpDir); err != nil || !info.IsDir() {
		t.Error("directory should have been created")
	}
}

func TestFileMkdirAll(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "seda_test_mkdir_all", "nested", "dir")
	defer os.RemoveAll(filepath.Join(os.TempDir(), "seda_test_mkdir_all"))

	input := `File.mkdir_all("` + tmpDir + `")`
	result := testEval(input)

	if result != object.NULL {
		t.Fatalf("mkdir_all should return NULL. got=%T (%+v)", result, result)
	}

	// Verify nested directories were created
	if info, err := os.Stat(tmpDir); err != nil || !info.IsDir() {
		t.Error("nested directories should have been created")
	}
}

func TestFileRemoveDir(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "seda_test_remove_dir")
	os.Mkdir(tmpDir, 0755)

	input := `File.remove_dir("` + tmpDir + `")`
	result := testEval(input)

	if result != object.NULL {
		t.Fatalf("remove_dir should return NULL. got=%T (%+v)", result, result)
	}

	// Verify directory was removed
	if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
		t.Error("directory should have been removed")
	}
}

func TestFileErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			`File.read()`,
			"wrong number of arguments for File.read. got=0, want=1",
		},
		{
			`File.read(123)`,
			"argument to File.read must be STRING, got NUMBER",
		},
		{
			`File.write("path")`,
			"wrong number of arguments for File.write. got=1, want=2",
		},
		{
			`File.write(123, "content")`,
			"first argument to File.write must be STRING, got NUMBER",
		},
		{
			`File.write("path", 123)`,
			"second argument to File.write must be STRING, got NUMBER",
		},
		{
			`File.delete()`,
			"wrong number of arguments for File.delete. got=0, want=1",
		},
		{
			`File.exists(123)`,
			"argument to File.exists must be STRING, got NUMBER",
		},
	}

	for _, tt := range tests {
		result := testEval(tt.input)

		errObj, ok := result.(*object.Error)
		if !ok {
			t.Errorf("Expected error for %s, got %T", tt.input, result)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("for %s:\nexpected error: %q\ngot error: %q",
				tt.input, tt.expectedMessage, errObj.Message)
		}
	}
}

func TestFileReadNonexistent(t *testing.T) {
	input := `
		var content, err = File.read("/nonexistent/file/path.txt")
		!isNull(err)
	`
	result := testEval(input)
	testBooleanObject(t, result, true)
}

func TestFileComplexOperations(t *testing.T) {
	// Create a temp directory for this test
	tmpDir := filepath.Join(os.TempDir(), "seda_test_complex")
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test.txt")

	// Write a file
	writeInput := `File.write("` + tmpFile + `", "Initial content")`
	testEval(writeInput)

	// Check exists
	existsInput := `File.exists("` + tmpFile + `")`
	existsResult := testEval(existsInput)
	testBooleanObject(t, existsResult, true)

	// Append to it
	appendInput := `File.append("` + tmpFile + `", "Appended line")`
	testEval(appendInput)

	// Get its size
	sizeInput := `
var size, err = File.size("` + tmpFile + `")
size
`
	sizeResult := testEval(sizeInput)

	num, ok := sizeResult.(*object.Number)
	if !ok {
		t.Fatalf("size is not Number. got=%T", sizeResult)
	}
	if num.Value <= 0 {
		t.Errorf("size should be > 0, got %f", num.Value)
	}

	// Verify it's a file
	isFileInput := `File.is_file("` + tmpFile + `")`
	isFileResult := testEval(isFileInput)
	testBooleanObject(t, isFileResult, true)
}
