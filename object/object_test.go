package object

import (
	"testing"

	"github.com/vpaulo/seda/ast"
)

// Test Number object
func TestNumberObject(t *testing.T) {
	tests := []struct {
		value    float64
		expected string
	}{
		{42.0, "42"},
		{3.14159, "3.14159"},
		{0.0, "0"},
		{-10.5, "-10.5"},
		{1000000.0, "1000000"},
	}

	for _, tt := range tests {
		num := &Number{Value: tt.value}

		if num.Type() != NUMBER_OBJ {
			t.Errorf("num.Type() = %q, want %q", num.Type(), NUMBER_OBJ)
		}

		if num.Inspect() != tt.expected {
			t.Errorf("num.Inspect() = %q, want %q", num.Inspect(), tt.expected)
		}

		if num.String() != tt.expected {
			t.Errorf("num.String() = %q, want %q", num.String(), tt.expected)
		}
	}
}

func TestNumberWithProperties(t *testing.T) {
	num := &Number{Value: 42, Properties: make(map[string]Object)}
	num.Properties["custom"] = &String{Value: "test"}

	if prop, ok := num.Properties["custom"]; ok {
		if prop.(*String).Value != "test" {
			t.Errorf("property value = %q, want %q", prop.(*String).Value, "test")
		}
	} else {
		t.Error("expected custom property to exist")
	}
}

// Test String object
func TestStringObject(t *testing.T) {
	tests := []struct {
		value           string
		expectedInspect string
		expectedString  string
	}{
		{"hello", `"hello"`, "hello"},
		{"world", `"world"`, "world"},
		{"", `""`, ""},
		{"with\nnewline", "\"with\nnewline\"", "with\nnewline"},
	}

	for _, tt := range tests {
		str := &String{Value: tt.value}

		if str.Type() != STRING_OBJ {
			t.Errorf("str.Type() = %q, want %q", str.Type(), STRING_OBJ)
		}

		if str.Inspect() != tt.expectedInspect {
			t.Errorf("str.Inspect() = %q, want %q", str.Inspect(), tt.expectedInspect)
		}

		if str.String() != tt.expectedString {
			t.Errorf("str.String() = %q, want %q", str.String(), tt.expectedString)
		}
	}
}

func TestStringWithProperties(t *testing.T) {
	str := &String{Value: "hello", Properties: make(map[string]Object)}
	str.Properties["reverse"] = &Builtin{Fn: func(args ...Object) Object {
		return &String{Value: "olleh"}
	}}

	if prop, ok := str.Properties["reverse"]; ok {
		if prop.Type() != BUILTIN_OBJ {
			t.Errorf("property type = %q, want %q", prop.Type(), BUILTIN_OBJ)
		}
	} else {
		t.Error("expected reverse property to exist")
	}
}

// Test Boolean object
func TestBooleanObject(t *testing.T) {
	tests := []struct {
		value    bool
		expected string
	}{
		{true, "true"},
		{false, "false"},
	}

	for _, tt := range tests {
		b := &Boolean{Value: tt.value}

		if b.Type() != BOOLEAN_OBJ {
			t.Errorf("b.Type() = %q, want %q", b.Type(), BOOLEAN_OBJ)
		}

		if b.Inspect() != tt.expected {
			t.Errorf("b.Inspect() = %q, want %q", b.Inspect(), tt.expected)
		}

		if b.String() != tt.expected {
			t.Errorf("b.String() = %q, want %q", b.String(), tt.expected)
		}
	}
}

func TestBooleanGlobalConstants(t *testing.T) {
	if TRUE.Value != true {
		t.Error("TRUE should have value true")
	}
	if FALSE.Value != false {
		t.Error("FALSE should have value false")
	}
	if TRUE.Type() != BOOLEAN_OBJ {
		t.Errorf("TRUE.Type() = %q, want %q", TRUE.Type(), BOOLEAN_OBJ)
	}
}

// Test Array object
func TestArrayObject(t *testing.T) {
	tests := []struct {
		name     string
		elements []Object
		expected string
	}{
		{
			"empty array",
			[]Object{},
			"[]",
		},
		{
			"single number",
			[]Object{&Number{Value: 42}},
			"[42]",
		},
		{
			"multiple numbers",
			[]Object{&Number{Value: 1}, &Number{Value: 2}, &Number{Value: 3}},
			"[1, 2, 3]",
		},
		{
			"mixed types",
			[]Object{&Number{Value: 1}, &String{Value: "hello"}, TRUE},
			`[1, "hello", true]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := &Array{Elements: tt.elements}

			if arr.Type() != ARRAY_OBJ {
				t.Errorf("arr.Type() = %q, want %q", arr.Type(), ARRAY_OBJ)
			}

			if arr.Inspect() != tt.expected {
				t.Errorf("arr.Inspect() = %q, want %q", arr.Inspect(), tt.expected)
			}

			if arr.String() != tt.expected {
				t.Errorf("arr.String() = %q, want %q", arr.String(), tt.expected)
			}
		})
	}
}

// Test Range object
func TestRangeObject(t *testing.T) {
	tests := []struct {
		name      string
		start     int
		end       int
		inclusive bool
		expected  string
	}{
		{"exclusive range", 1, 10, false, "1..10"},
		{"inclusive range", 1, 10, true, "1...10"},
		{"zero start", 0, 5, false, "0..5"},
		{"negative range", -5, 5, true, "-5...5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Range{Start: tt.start, End: tt.end, Inclusive: tt.inclusive}

			if r.Type() != RANGE_OBJ {
				t.Errorf("r.Type() = %q, want %q", r.Type(), RANGE_OBJ)
			}

			if r.Inspect() != tt.expected {
				t.Errorf("r.Inspect() = %q, want %q", r.Inspect(), tt.expected)
			}

			if r.String() != tt.expected {
				t.Errorf("r.String() = %q, want %q", r.String(), tt.expected)
			}
		})
	}
}

// Test Map object
func TestMapObject(t *testing.T) {
	tests := []struct {
		name     string
		pairs    map[string]MapPair
		expected string
	}{
		{
			"empty map",
			map[string]MapPair{},
			"{}",
		},
		{
			"single pair",
			map[string]MapPair{
				"key1": {Key: &String{Value: "key1"}, Value: &Number{Value: 42}},
			},
			`{"key1": 42}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{Pairs: tt.pairs}

			if m.Type() != MAP_OBJ {
				t.Errorf("m.Type() = %q, want %q", m.Type(), MAP_OBJ)
			}

			if m.String() != tt.expected {
				t.Errorf("m.String() = %q, want %q", m.String(), tt.expected)
			}
		})
	}
}

// Test Null object
func TestNullObject(t *testing.T) {
	null := &Null{}

	if null.Type() != NULL_OBJ {
		t.Errorf("null.Type() = %q, want %q", null.Type(), NULL_OBJ)
	}

	if null.Inspect() != "null" {
		t.Errorf("null.Inspect() = %q, want %q", null.Inspect(), "null")
	}

	if null.String() != "null" {
		t.Errorf("null.String() = %q, want %q", null.String(), "null")
	}

	// Test global constant
	if NULL.Type() != NULL_OBJ {
		t.Errorf("NULL.Type() = %q, want %q", NULL.Type(), NULL_OBJ)
	}
}

// Test Error object
func TestErrorObject(t *testing.T) {
	tests := []struct {
		message         string
		expectedInspect string
		expectedString  string
	}{
		{"simple error", "ERROR: simple error", "simple error"},
		{"type mismatch", "ERROR: type mismatch", "type mismatch"},
	}

	for _, tt := range tests {
		err := &Error{Message: tt.message}

		if err.Type() != ERROR_OBJ {
			t.Errorf("err.Type() = %q, want %q", err.Type(), ERROR_OBJ)
		}

		if err.Inspect() != tt.expectedInspect {
			t.Errorf("err.Inspect() = %q, want %q", err.Inspect(), tt.expectedInspect)
		}

		if err.String() != tt.expectedString {
			t.Errorf("err.String() = %q, want %q", err.String(), tt.expectedString)
		}
	}
}

func TestNewError(t *testing.T) {
	err := NewError("error: %s %d", "test", 42)

	if err.Message != "error: test 42" {
		t.Errorf("err.Message = %q, want %q", err.Message, "error: test 42")
	}

	if err.Type() != ERROR_OBJ {
		t.Errorf("err.Type() = %q, want %q", err.Type(), ERROR_OBJ)
	}
}

// Test ReturnValue object
func TestReturnValueObject(t *testing.T) {
	innerValue := &Number{Value: 42}
	rv := &ReturnValue{Value: innerValue}

	if rv.Type() != RETURN_VALUE_OBJ {
		t.Errorf("rv.Type() = %q, want %q", rv.Type(), RETURN_VALUE_OBJ)
	}

	if rv.Inspect() != "42" {
		t.Errorf("rv.Inspect() = %q, want %q", rv.Inspect(), "42")
	}

	if rv.String() != "42" {
		t.Errorf("rv.String() = %q, want %q", rv.String(), "42")
	}
}

// Test Break object
func TestBreakObject(t *testing.T) {
	b := &Break{}

	if b.Type() != BREAK_OBJ {
		t.Errorf("b.Type() = %q, want %q", b.Type(), BREAK_OBJ)
	}

	if b.Inspect() != "break" {
		t.Errorf("b.Inspect() = %q, want %q", b.Inspect(), "break")
	}

	if b.String() != "break" {
		t.Errorf("b.String() = %q, want %q", b.String(), "break")
	}
}

// Test IsTruthy function
func TestIsTruthy(t *testing.T) {
	tests := []struct {
		name     string
		obj      Object
		expected bool
	}{
		{"null is falsy", NULL, false},
		{"true is truthy", TRUE, true},
		{"false is falsy", FALSE, false},
		{"number is truthy", &Number{Value: 0}, true},
		{"string is truthy", &String{Value: ""}, true},
		{"array is truthy", &Array{Elements: []Object{}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTruthy(tt.obj)
			if result != tt.expected {
				t.Errorf("IsTruthy(%v) = %v, want %v", tt.obj, result, tt.expected)
			}
		})
	}
}

// Test IsError function
func TestIsError(t *testing.T) {
	tests := []struct {
		name     string
		obj      Object
		expected bool
	}{
		{"error object", &Error{Message: "test"}, true},
		{"nil object", nil, false},
		{"number object", &Number{Value: 42}, false},
		{"string object", &String{Value: "test"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsError(tt.obj)
			if result != tt.expected {
				t.Errorf("IsError(%v) = %v, want %v", tt.obj, result, tt.expected)
			}
		})
	}
}

// Test Function object
func TestFunctionObject(t *testing.T) {
	params := []*ast.Parameter{
		{Name: &ast.Identifier{Value: "x"}},
		{Name: &ast.Identifier{Value: "y"}},
	}
	body := &ast.BlockStatement{
		Statements: []ast.Statement{},
	}
	env := NewEnvironment()

	fn := &Function{
		Parameters: params,
		Body:       body,
		Env:        env,
	}

	if fn.Type() != FUNCTION_OBJ {
		t.Errorf("fn.Type() = %q, want %q", fn.Type(), FUNCTION_OBJ)
	}

	inspect := fn.Inspect()
	if len(inspect) == 0 {
		t.Error("fn.Inspect() should not be empty")
	}
}

// Test Builtin object
func TestBuiltinObject(t *testing.T) {
	builtinFn := func(args ...Object) Object {
		return &Number{Value: 42}
	}

	builtin := &Builtin{Fn: builtinFn}

	if builtin.Type() != BUILTIN_OBJ {
		t.Errorf("builtin.Type() = %q, want %q", builtin.Type(), BUILTIN_OBJ)
	}

	if builtin.Inspect() != "builtin function" {
		t.Errorf("builtin.Inspect() = %q, want %q", builtin.Inspect(), "builtin function")
	}

	result := builtin.Fn()
	if result.(*Number).Value != 42 {
		t.Errorf("builtin.Fn() = %v, want 42", result)
	}
}

// Test Module object
func TestModuleObject(t *testing.T) {
	env := NewEnvironment()
	mod := &Module{
		Name:        "TestModule",
		Environment: env,
	}

	if mod.Type() != MODULE_OBJ {
		t.Errorf("mod.Type() = %q, want %q", mod.Type(), MODULE_OBJ)
	}

	if mod.Inspect() != "module TestModule" {
		t.Errorf("mod.Inspect() = %q, want %q", mod.Inspect(), "module TestModule")
	}

	if mod.String() != "module TestModule" {
		t.Errorf("mod.String() = %q, want %q", mod.String(), "module TestModule")
	}
}

// Test TypeAlias object
func TestTypeAliasObject(t *testing.T) {
	typeAnnotation := &ast.TypeAnnotation{
		Name: "number",
	}

	typeAlias := &TypeAlias{
		Name:           "MyNumber",
		TypeAnnotation: typeAnnotation,
	}

	if typeAlias.Type() != TYPE_ALIAS_OBJ {
		t.Errorf("typeAlias.Type() = %q, want %q", typeAlias.Type(), TYPE_ALIAS_OBJ)
	}

	if typeAlias.Inspect() != "type MyNumber" {
		t.Errorf("typeAlias.Inspect() = %q, want %q", typeAlias.Inspect(), "type MyNumber")
	}

	if typeAlias.String() != "type MyNumber" {
		t.Errorf("typeAlias.String() = %q, want %q", typeAlias.String(), "type MyNumber")
	}
}

// Test TestResult object
func TestTestResultObject(t *testing.T) {
	tests := []struct {
		name           string
		testResult     *TestResult
		expectedStatus string
	}{
		{
			"all passed",
			&TestResult{Passed: 3, Failed: 0, Failures: []string{}, Label: ""},
			"PASSED",
		},
		{
			"some failed",
			&TestResult{Passed: 2, Failed: 1, Failures: []string{"assertion failed"}, Label: ""},
			"FAILED",
		},
		{
			"with label",
			&TestResult{Passed: 5, Failed: 0, Failures: []string{}, Label: "my test"},
			"PASSED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.testResult.Type() != TEST_RESULT_OBJ {
				t.Errorf("testResult.Type() = %q, want %q", tt.testResult.Type(), TEST_RESULT_OBJ)
			}

			inspect := tt.testResult.Inspect()
			if len(inspect) == 0 {
				t.Error("testResult.Inspect() should not be empty")
			}

			// Check if status appears in output
			if tt.expectedStatus == "PASSED" && tt.testResult.Failed == 0 {
				// Verify passed count
				if tt.testResult.Passed == 0 {
					t.Error("expected passed count > 0")
				}
			}

			if tt.expectedStatus == "FAILED" && tt.testResult.Failed > 0 {
				// Verify failures are included
				if len(tt.testResult.Failures) == 0 {
					t.Error("expected failures to be recorded")
				}
			}
		})
	}
}

// Test Environment
func TestEnvironment(t *testing.T) {
	env := NewEnvironment()

	// Test Set and Get
	num := &Number{Value: 42}
	env.Set("x", num)

	val, ok := env.Get("x")
	if !ok {
		t.Error("expected to find 'x' in environment")
	}

	if val.(*Number).Value != 42 {
		t.Errorf("val = %v, want 42", val)
	}
}

func TestEnvironmentOuter(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("x", &Number{Value: 10})

	inner := NewEnclosedEnvironment(outer)
	inner.Set("y", &Number{Value: 20})

	// Inner can access outer
	val, ok := inner.Get("x")
	if !ok {
		t.Error("expected inner to access outer variable 'x'")
	}
	if val.(*Number).Value != 10 {
		t.Errorf("val = %v, want 10", val)
	}

	// Inner has its own variables
	val, ok = inner.Get("y")
	if !ok {
		t.Error("expected to find 'y' in inner environment")
	}
	if val.(*Number).Value != 20 {
		t.Errorf("val = %v, want 20", val)
	}

	// Outer doesn't have inner's variables
	_, ok = outer.Get("y")
	if ok {
		t.Error("expected outer to not have inner variable 'y'")
	}
}

func TestEnvironmentConstants(t *testing.T) {
	env := NewEnvironment()

	// Set a constant
	env.SetConstant("PI", &Number{Value: 3.14159})

	// Check if it's a constant
	if !env.IsConstant("PI") {
		t.Error("expected 'PI' to be a constant")
	}

	// Try to update constant
	result := env.Update("PI", &Number{Value: 3.0})
	if !IsError(result) {
		t.Error("expected error when reassigning constant")
	}

	// Set a regular variable
	env.Set("x", &Number{Value: 10})

	// Should not be a constant
	if env.IsConstant("x") {
		t.Error("expected 'x' to not be a constant")
	}

	// Should allow update
	result = env.Update("x", &Number{Value: 20})
	if IsError(result) {
		t.Errorf("unexpected error when updating variable: %v", result)
	}

	val, _ := env.Get("x")
	if val.(*Number).Value != 20 {
		t.Errorf("val = %v, want 20", val)
	}
}

func TestEnvironmentUpdate(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("x", &Number{Value: 10})

	inner := NewEnclosedEnvironment(outer)

	// Update outer variable from inner scope
	inner.Update("x", &Number{Value: 20})

	val, _ := outer.Get("x")
	if val.(*Number).Value != 20 {
		t.Errorf("val = %v, want 20 (outer variable should be updated)", val)
	}

	// Update non-existent variable creates it in current scope
	inner.Update("y", &Number{Value: 30})

	val, ok := inner.Get("y")
	if !ok {
		t.Error("expected 'y' to be created in inner scope")
	}
	if val.(*Number).Value != 30 {
		t.Errorf("val = %v, want 30", val)
	}

	// Outer shouldn't have the new variable
	_, ok = outer.Get("y")
	if ok {
		t.Error("expected outer to not have 'y'")
	}
}

func TestEnvironmentIsInWhereBlockTest(t *testing.T) {
	outer := NewEnvironment()
	outer.InWhereBlockTest = false

	inner := NewEnclosedEnvironment(outer)
	inner.InWhereBlockTest = true

	if !inner.IsInWhereBlockTest() {
		t.Error("expected IsInWhereBlockTest to return true for inner")
	}

	innerInner := NewEnclosedEnvironment(inner)
	if !innerInner.IsInWhereBlockTest() {
		t.Error("expected IsInWhereBlockTest to propagate to nested environments")
	}

	standalone := NewEnvironment()
	if standalone.IsInWhereBlockTest() {
		t.Error("expected IsInWhereBlockTest to return false for standalone environment")
	}
}

func TestEnvironmentGetStore(t *testing.T) {
	env := NewEnvironment()
	env.Set("x", &Number{Value: 42})
	env.Set("y", &String{Value: "hello"})

	store := env.GetStore()
	if len(store) != 2 {
		t.Errorf("store length = %d, want 2", len(store))
	}

	if _, ok := store["x"]; !ok {
		t.Error("expected store to contain 'x'")
	}

	if _, ok := store["y"]; !ok {
		t.Error("expected store to contain 'y'")
	}
}
