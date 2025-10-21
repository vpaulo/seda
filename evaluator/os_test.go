package evaluator

import (
	"os"
	"runtime"
	"testing"

	"github.com/vpaulo/seda/object"
)

// OS Module Tests

func TestOSGetEnv(t *testing.T) {
	// Set a test environment variable
	os.Setenv("TEST_VAR_SEDA", "test_value")
	defer os.Unsetenv("TEST_VAR_SEDA")

	input := `OS.getenv("TEST_VAR_SEDA")`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	if str.Value != "test_value" {
		t.Errorf("wrong value. expected=%s, got=%s", "test_value", str.Value)
	}
}

func TestOSSetEnv(t *testing.T) {
	input := `
		var _ = OS.setenv("TEST_SEDA_SET", "new_value")
		OS.getenv("TEST_SEDA_SET")
	`
	result := testEval(input)
	defer os.Unsetenv("TEST_SEDA_SET")

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	if str.Value != "new_value" {
		t.Errorf("wrong value. expected=%s, got=%s", "new_value", str.Value)
	}
}

func TestOSEnv(t *testing.T) {
	// Set a known environment variable
	os.Setenv("TEST_ENV_MAP", "map_value")
	defer os.Unsetenv("TEST_ENV_MAP")

	input := `
		var env_map = OS.env()
		env_map.TEST_ENV_MAP
	`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	if str.Value != "map_value" {
		t.Errorf("wrong value. expected=%s, got=%s", "map_value", str.Value)
	}
}

func TestOSArgs(t *testing.T) {
	// Save original args
	originalArgs := command_line_args
	defer func() { command_line_args = originalArgs }()

	// Set test args
	command_line_args = []string{"program", "arg1", "arg2"}

	input := `OS.args()`
	result := testEval(input)

	arr, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("result is not Array. got=%T (%+v)", result, result)
	}

	if len(arr.Elements) != 3 {
		t.Fatalf("array has wrong length. expected=3, got=%d", len(arr.Elements))
	}

	expected := []string{"program", "arg1", "arg2"}
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

func TestOSPID(t *testing.T) {
	input := `OS.pid()`
	result := testEval(input)

	num, ok := result.(*object.Number)
	if !ok {
		t.Fatalf("result is not Number. got=%T (%+v)", result, result)
	}

	expectedPID := float64(os.Getpid())
	if num.Value != expectedPID {
		t.Errorf("wrong PID. expected=%f, got=%f", expectedPID, num.Value)
	}
}

func TestOSPlatform(t *testing.T) {
	input := `OS.platform()`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	if str.Value != runtime.GOOS {
		t.Errorf("wrong platform. expected=%s, got=%s", runtime.GOOS, str.Value)
	}
}

func TestOSArch(t *testing.T) {
	input := `OS.arch()`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	if str.Value != runtime.GOARCH {
		t.Errorf("wrong arch. expected=%s, got=%s", runtime.GOARCH, str.Value)
	}
}

func TestOSHostname(t *testing.T) {
	input := `OS.hostname()`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	expectedHostname, err := os.Hostname()
	if err != nil {
		t.Fatalf("failed to get hostname: %v", err)
	}

	if str.Value != expectedHostname {
		t.Errorf("wrong hostname. expected=%s, got=%s", expectedHostname, str.Value)
	}
}

func TestOSHomeDir(t *testing.T) {
	input := `OS.home_dir()`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	expectedHome, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	if str.Value != expectedHome {
		t.Errorf("wrong home dir. expected=%s, got=%s", expectedHome, str.Value)
	}
}

func TestOSTempDir(t *testing.T) {
	input := `OS.temp_dir()`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	expectedTemp := os.TempDir()
	if str.Value != expectedTemp {
		t.Errorf("wrong temp dir. expected=%s, got=%s", expectedTemp, str.Value)
	}
}

func TestOSCwd(t *testing.T) {
	input := `OS.cwd()`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	expectedCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	if str.Value != expectedCwd {
		t.Errorf("wrong cwd. expected=%s, got=%s", expectedCwd, str.Value)
	}
}

func TestOSErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			`OS.getenv()`,
			"wrong number of arguments for OS.getenv. got=0, want=1",
		},
		{
			`OS.getenv(123)`,
			"argument to OS.getenv must be STRING, got NUMBER",
		},
		{
			`OS.setenv("VAR")`,
			"wrong number of arguments for OS.setenv. got=1, want=2",
		},
		{
			`OS.setenv(123, "value")`,
			"first argument to OS.setenv must be STRING, got NUMBER",
		},
		{
			`OS.setenv("VAR", 123)`,
			"second argument to OS.setenv must be STRING, got NUMBER",
		},
		{
			`OS.env(123)`,
			"wrong number of arguments for OS.env. got=1, want=0",
		},
		{
			`OS.args(123)`,
			"wrong number of arguments for OS.args. got=1, want=0",
		},
		{
			`OS.pid(123)`,
			"wrong number of arguments for OS.pid. got=1, want=0",
		},
		{
			`OS.platform(123)`,
			"wrong number of arguments for OS.platform. got=1, want=0",
		},
		{
			`OS.arch(123)`,
			"wrong number of arguments for OS.arch. got=1, want=0",
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

func TestOSComplexOperations(t *testing.T) {
	// Set multiple environment variables and retrieve them
	input := `
		var _ = OS.setenv("VAR1", "value1")
		var _ = OS.setenv("VAR2", "value2")
		var _ = OS.setenv("VAR3", "value3")

		var v1 = OS.getenv("VAR1")
		var v2 = OS.getenv("VAR2")
		var v3 = OS.getenv("VAR3")

		v1 + "," + v2 + "," + v3
	`
	result := testEval(input)
	defer func() {
		os.Unsetenv("VAR1")
		os.Unsetenv("VAR2")
		os.Unsetenv("VAR3")
	}()

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	expected := "value1,value2,value3"
	if str.Value != expected {
		t.Errorf("wrong result. expected=%s, got=%s", expected, str.Value)
	}
}

func TestOSSystemInfo(t *testing.T) {
	// Test combining system info
	input := `
		var platform = OS.platform()
		var arch = OS.arch()
		var info = platform + "/" + arch
		info
	`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	expected := runtime.GOOS + "/" + runtime.GOARCH
	if str.Value != expected {
		t.Errorf("wrong result. expected=%s, got=%s", expected, str.Value)
	}
}
