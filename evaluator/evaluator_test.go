package evaluator

import (
	"testing"

	"github.com/vpaulo/seda/lexer"
	"github.com/vpaulo/seda/object"
	"github.com/vpaulo/seda/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func TestEvalNumberExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a;", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected.(int))
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"var i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"var myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"var myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"var myArray = [1, 2, 3]; var i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

// Helper functions
func testNumberObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Number)
	if !ok {
		t.Errorf("object is not Number. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%g, want=%g", result.Value, expected)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int) bool {
	return testNumberObject(t, obj, float64(expected))
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != object.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

// Control flow tests

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true :: 10 end", 10},
		{"if false :: 10 end", nil},
		{"if 1 :: 10 end", 10},
		{"if 1 < 2 :: 10 end", 10},
		{"if 1 > 2 :: 10 end", nil},
		{"if 1 > 2 :: 10 else :: 20 end", 20},
		{"if 1 < 2 :: 10 else :: 20 end", 10},
		{
			`if false ::
				10
			else if true ::
				20
			else ::
				30
			end`,
			20,
		},
		{
			`if false ::
				10
			else if false ::
				20
			else ::
				30
			end`,
			30,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestForStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`var sum = 0
			for i in [1, 2, 3] ::
				sum + i
			end`,
			3, // Returns last evaluated expression in loop
		},
		{
			`var text = ""
			for char in "abc" ::
				char
			end`,
			"c", // Returns last character
		},
		{
			`for i, val in [10, 20, 30] ::
				i + val
			end`,
			32, // Last iteration: 2 + 30 = 32
		},
		{
			`for i in [1, 2, 3] ::
				i * 2
			end`,
			6, // Last iteration: 3 * 2 = 6
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, expected)
		case string:
			testStringObject(t, evaluated, expected)
		}
	}
}

func TestCaseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`case 1 ::
				1 => "one"
				2 => "two"
				_ => "other"
			end`,
			"one",
		},
		{
			`case 2 ::
				1 => "one"
				2 => "two"
				_ => "other"
			end`,
			"two",
		},
		{
			`case 3 ::
				1 => "one"
				2 => "two"
				_ => "other"
			end`,
			"other",
		},
		{
			`var x = 5
			case x ::
				3 => "three"
				5 => "five"
				_ => "unknown"
			end`,
			"five",
		},
		{
			`case "hello" ::
				"hi" => 1
				"hello" => 2
				"bye" => 3
				_ => 0
			end`,
			2,
		},
		{
			`case true ::
				true => "yes"
				false => "no"
				_ => "maybe"
			end`,
			"yes",
		},
		{
			`case 2 + 3 ::
				4 => "four"
				5 => "five"
				6 => "six"
				_ => "other"
			end`,
			"five",
		},
		{
			// Test case expression as variable assignment
			`var grade = "A"
			var result = case grade ::
				"A" => 4.0
				"B" => 3.0
				_ => 0.0
			end
			result`,
			4.0,
		},
		{
			// Test case expression with no match (should return NULL)
			`case "X" ::
				"A" => "excellent"
				"B" => "good"
			end`,
			nil, // NULL
		},
		{
			// Test nested case expressions
			`var category = "electronics"
			case category ::
				"electronics" => case "phone" ::
				                   "phone" => "Mobile"
				                   "laptop" => "Computer"
				                   _ => "Other"
				                 end
				"clothing" => "Fashion"
				_ => "Unknown"
			end`,
			"Mobile",
		},
	}

	for i, tt := range tests {
		obj := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testNumberObject(t, obj, float64(expected))
		case float64:
			testNumberObject(t, obj, expected)
		case string:
			testStringObject(t, obj, expected)
		case nil:
			testNullObject(t, obj)
		default:
			t.Errorf("test[%d]: unknown expected type: %T", i, expected)
		}
	}
}

func TestCaseStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`case 1 ::
				1 => "one"
				2 => "two"
				_ => "other"
			end`,
			"one",
		},
		{
			`case 2 ::
				1 => "one"
				2 => "two"
				_ => "other"
			end`,
			"two",
		},
		{
			`case 3 ::
				1 => "one"
				2 => "two"
				_ => "other"
			end`,
			"other",
		},
		{
			`var x = 5
			case x ::
				1 => 10
				5 => 50
				10 => 100
			end`,
			50,
		},
		{
			`case "hello" ::
				"hi" => 1
				"hello" => 2
				"bye" => 3
			end`,
			2,
		},
		{
			`case true ::
				false => "no"
				true => "yes"
			end`,
			"yes",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, expected)
		case string:
			testStringObject(t, evaluated, expected)
		}
	}
}

func TestNestedControlFlow(t *testing.T) {
	input := `
	for i in [1, 2, 3, 4, 5] ::
		if i % 2 == 0 ::
			i + 10
		else ::
			i - 10
		end
	end
	`

	evaluated := testEval(input)
	// Last iteration (i=5): 5 % 2 == 1 (odd), so 5 - 10 = -5
	testIntegerObject(t, evaluated, -5)
}

func TestComplexControlFlow(t *testing.T) {
	input := `
	var numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
	var lastEven = 0
	var lastOdd = 0

	for num in numbers ::
		case num % 2 ::
			0 => num
			1 => num + 100
		end
	end
	`

	evaluated := testEval(input)
	// Last iteration: num=10, 10 % 2 == 0, so return 10
	testIntegerObject(t, evaluated, 10)
}

// Helper function for string testing
func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
		return false
	}
	return true
}

// Function and method tests

func TestFunctionDefinitions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`fn add(x: number, y: number): number ::
				x + y
			end
			add`,
			"function object",
		},
		{
			`fn greet(name: string): string ::
				"Hello, " + name
			end
			greet("World")`,
			"Hello, World",
		},
		{
			`fn factorial(n: number): number ::
				if n <= 1 ::
					1
				else ::
					n * factorial(n - 1)
				end
			end
			factorial(5)`,
			120,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, expected)
		case string:
			if expected == "function object" {
				_, ok := evaluated.(*object.Function)
				if !ok {
					t.Errorf("object is not Function. got=%T (%+v)", evaluated, evaluated)
				}
			} else {
				testStringObject(t, evaluated, expected)
			}
		}
	}
}

func TestFunctionLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`var double = fn(x: number) :: x * 2 end
			double(5)`,
			10,
		},
		{
			`var add = fn(x: number, y: number) :: x + y end
			add(3, 4)`,
			7,
		},
		{
			`fn(x: number) :: x + 1 end`,
			"function object",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, expected)
		case string:
			if expected == "function object" {
				_, ok := evaluated.(*object.Function)
				if !ok {
					t.Errorf("object is not Function. got=%T (%+v)", evaluated, evaluated)
				}
			}
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`fn early_return(x: number): number ::
				if x > 5 ::
					return x * 2
				end
				x + 1
			end
			early_return(10)`,
			20,
		},
		{
			`fn early_return(x: number): number ::
				if x > 5 ::
					return x * 2
				end
				x + 1
			end
			early_return(3)`,
			4,
		},
		{
			`fn void_function() ::
				return
			end
			void_function()`,
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestSimpleFunctionCall(t *testing.T) {
	input := `
	fn add(x: number, y: number) ::
		x + y
	end
	add(2, 3)
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 5)
}

func TestClosures(t *testing.T) {
	// Simpler closure test first
	input := `
	fn outer() ::
		var x = 10
		fn inner() ::
			x
		end
		inner
	end
	var f = outer()
	f()
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10)
}

// Note: Method definitions and calls are deferred until struct support is implemented
// This would require:
// 1. Struct object implementation
// 2. Method dispatch system
// 3. Self binding in method calls
// 4. Type-based method storage

// Builtin function tests

func TestLengthMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"".length()`, 0},
		{`"hello".length()`, 5},
		{`"hello world".length()`, 11},
		{`[].length()`, 0},
		{`[1, 2, 3].length()`, 3},
		{`[1, 2, 3, 4, 5].length()`, 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, expected)
		}
	}
}

func TestStringMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"hello".upper()`, "HELLO"},
		{`"Hello World".upper()`, "HELLO WORLD"},
		{`"HELLO".lower()`, "hello"},
		{`"Hello World".lower()`, "hello world"},
		{`"hello".substr(1, 3)`, "ell"},
		{`"hello world".substr(0, 5)`, "hello"},
		{`"hello world".substr(6, 5)`, "world"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case string:
			testStringObject(t, evaluated, expected)
		}
	}
}

func TestArrayMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`[1, 2, 3].first()`, 1},
		{`[].first()`, nil},
		{`[1, 2, 3].last()`, 3},
		{`[].last()`, nil},
		{`[1, 2, 3].pop()`, 3},
		{`[].pop()`, nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestPushMethod(t *testing.T) {
	input := `
	var arr = [1, 2, 3]
	var newArr = arr.push(4)
	newArr.length()
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 4)
}

func TestMethodErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			`42.length()`,
			"method 'length' not found on Number",
		},
		{
			`"hello".length("arg")`,
			"wrong number of arguments for String.length. got=1, want=0",
		},
		{
			`42.upper()`,
			"method 'upper' not found on Number",
		},
		{
			`"hello".substr("world", 1)`,
			"first argument to String.substr must be NUMBER, got STRING",
		},
		{
			`42.push(1)`,
			"method 'push' not found on Number",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestPrintBuiltins(t *testing.T) {
	// Note: print and println write to stdout, so we test they don't error
	// In a real implementation, you might want to capture stdout for testing
	tests := []struct {
		input string
	}{
		{`print("hello")`},
		{`print("hello", "world")`},
		{`println("hello")`},
		{`println("hello", "world")`},
		{`print(42)`},
		{`println(true)`},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		// print/println should return NULL
		testNullObject(t, evaluated)
	}
}

func TestObjectOrientedMethodChaining(t *testing.T) {
	// Test that demonstrates the object-oriented nature of the language
	input := `
	var text = "hello world"
	var processedText = text.upper().substr(0, 5)
	processedText.length()
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 5)
}

func TestComplexObjectMethodUsage(t *testing.T) {
	// Test complex usage combining arrays and strings with methods
	input := `
	var words = ["hello", "beautiful", "world"]
	var firstWord = words.first()
	var upperFirst = firstWord.upper()
	var result = words.push(upperFirst)
	result.length()
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 4)
}

func TestObjectMethodsWithVariables(t *testing.T) {
	// Test that methods work with variables containing objects
	input := `
	var numbers = [1, 2, 3, 4, 5]
	var lastNum = numbers.last()
	var strNum = lastNum.to_string()
	strNum.length()
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 1) // "5" has length 1
}

// Testing Framework Tests

func TestCheckStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		failed   int
	}{
		{
			`var x = 5
			 check "basic test" ::
			   x is 5
			 end`,
			"Test 'basic test' PASSED: 1/1 assertions passed",
			0,
		},
		{
			`var x = 5
			 check "failing test" ::
			   x is 10
			 end`,
			"Test 'failing test' FAILED: 0/1 assertions passed",
			1,
		},
		{
			`var x = 5
			 var y = "hello"
			 check "mixed test" ::
			   x is 5
			   y is "hello"
			   x isA "Number"
			   y isA "String"
			 end`,
			"Test 'mixed test' PASSED: 4/4 assertions passed",
			0,
		},
		{
			`var arr = [1, 2, 3]
			 check "array test" ::
			   arr contains 2
			   arr contains 5
			 end`,
			"Test 'array test' FAILED: 1/2 assertions passed",
			1,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if testResult, ok := evaluated.(*object.TestResult); ok {
			if testResult.Failed != tt.failed {
				t.Errorf("Expected %d failed assertions, got %d", tt.failed, testResult.Failed)
			}

			resultStr := testResult.String()
			if !contains(resultStr, tt.expected) {
				t.Errorf("Expected result to contain %q, got %q", tt.expected, resultStr)
			}
		} else {
			t.Errorf("Expected TestResult object, got %T", evaluated)
		}
	}
}

func TestAssertionOperators(t *testing.T) {
	tests := []struct {
		input  string
		passed int
		failed int
	}{
		{
			`var x = 42
			 check "is operator" ::
			   x is 42
			   x is 43
			 end`,
			1, 1,
		},
		{
			`var str = "hello"
			 var num = 123
			 check "isA operator" ::
			   str isA "String"
			   num isA "Number"
			   str isA "Number"
			 end`,
			2, 1,
		},
		{
			`var arr = [1, 2, 3, 4]
			 var text = "hello world"
			 check "contains operator" ::
			   arr contains 3
			   arr contains 7
			   text contains "world"
			   text contains "xyz"
			 end`,
			2, 2,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if testResult, ok := evaluated.(*object.TestResult); ok {
			if testResult.Passed != tt.passed {
				t.Errorf("Expected %d passed assertions, got %d", tt.passed, testResult.Passed)
			}
			if testResult.Failed != tt.failed {
				t.Errorf("Expected %d failed assertions, got %d", tt.failed, testResult.Failed)
			}
		} else {
			t.Errorf("Expected TestResult object, got %T", evaluated)
		}
	}
}

func TestWhereBlocks(t *testing.T) {
	input := `
	fn add(a, b) ::
	  return a + b
	where ::
	  result is 8
	  arg0 is 3
	  arg1 is 5
	end

	add(3, 5)
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 8)
}

// Map Literal Tests

func TestMapLiterals(t *testing.T) {
	input := `var two = "two"
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Map)
	if !ok {
		t.Fatalf("Eval didn't return Map. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
		"4":     4,
		"true":  5,
		"false": 6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Map has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs. key=%q", expectedKey)
			continue
		}

		testIntegerObject(t, pair.Value, int(expectedValue))
	}
}

func TestMapIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`var key = "foo"
			 {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

// Assignment Expression Tests

func TestAssignmentExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{
			`var x = 5
			 x = 10
			 x`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// Comment Tests

func TestCommentsIgnored(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{
			`# This is a comment
			 var x = 5  # Another comment
			 x  # Final comment`,
			5,
		},
		{
			`# Multiple
			 # comment
			 # lines
			 var result = 42
			 result`,
			42,
		},
		{
			`fn add(a, b) ::  # Function comment
			   # Comment inside function
			   return a + b  # Return comment
			 end
			 add(3, 7)`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// Error Handling Tests

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true",
			"type mismatch: NUMBER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: NUMBER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) :: true + false end",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

// Edge Cases and Integration Tests

func TestDivisionByZero(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 / 0", "division by zero"},
		{"10 % 0", "division by zero"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("Expected error object, got %T", evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("Expected error message %q, got %q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`var a = [1, 2, 3]
			 var b = {"x": 10, "y": 20}
			 a[0] + b["x"]`,
			11,
		},
		{
			`var counter = fn(x) ::
			   fn(y) :: return x + y end
			 end
			 var addTwo = counter(2)
			 addTwo(3)`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, expected)
		default:
			t.Errorf("Unsupported expected type: %T", expected)
		}
	}
}

// Helper function for string containment check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				findInString(s, substr))))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Logical Operators Tests

func TestLogicalAndOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true && true", true},
		{"true && false", false},
		{"false && true", false},
		{"false && false", false},
		{"(1 < 2) && (3 < 4)", true},
		{"(1 < 2) && (3 > 4)", false},
		{"(1 > 2) && (3 < 4)", false},
		{"(5 == 5) && (10 == 10)", true},
		{"(5 == 5) && (10 != 10)", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestLogicalOrOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true || true", true},
		{"true || false", true},
		{"false || true", true},
		{"false || false", false},
		{"(1 < 2) || (3 < 4)", true},
		{"(1 < 2) || (3 > 4)", true},
		{"(1 > 2) || (3 < 4)", true},
		{"(1 > 2) || (3 > 4)", false},
		{"(5 != 5) || (10 == 10)", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestLogicalOperatorShortCircuit(t *testing.T) {
	// Test that && short-circuits on false
	input1 := `
	var x = 0
	if false && (x = 1) ::
		x
	end
	x
	`
	evaluated := testEval(input1)
	testIntegerObject(t, evaluated, 0) // x should still be 0

	// Test that || short-circuits on true
	input2 := `
	var y = 0
	if true || (y = 1) ::
		y
	else ::
		y
	end
	`
	evaluated = testEval(input2)
	testIntegerObject(t, evaluated, 0) // y should still be 0
}

func TestComplexLogicalExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true && true && true", true},
		{"true && true && false", false},
		{"false || false || true", true},
		{"false || false || false", false},
		{"(true && false) || true", true},
		{"true && (false || true)", true},
		{"(1 < 2) && (3 < 4) && (5 < 6)", true},
		{"(1 < 2) || (3 > 4) && (5 < 6)", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

// Range Expression Tests

func TestRangeExpressions(t *testing.T) {
	tests := []struct {
		input       string
		start       int
		end         int
		inclusive   bool
		description string
	}{
		{"1..10", 1, 10, false, "exclusive range"},
		{"1...10", 1, 10, true, "inclusive range"},
		{"0..5", 0, 5, false, "zero start exclusive"},
		{"0...5", 0, 5, true, "zero start inclusive"},
		{"-5..5", -5, 5, false, "negative to positive exclusive"},
		{"-5...5", -5, 5, true, "negative to positive inclusive"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		rangeObj, ok := evaluated.(*object.Range)
		if !ok {
			t.Errorf("%s: object is not Range. got=%T (%+v)", tt.description, evaluated, evaluated)
			continue
		}

		if rangeObj.Start != tt.start {
			t.Errorf("%s: wrong start. got=%d, want=%d", tt.description, rangeObj.Start, tt.start)
		}

		if rangeObj.End != tt.end {
			t.Errorf("%s: wrong end. got=%d, want=%d", tt.description, rangeObj.End, tt.end)
		}

		if rangeObj.Inclusive != tt.inclusive {
			t.Errorf("%s: wrong inclusive. got=%t, want=%t", tt.description, rangeObj.Inclusive, tt.inclusive)
		}
	}
}

func TestRangeInForLoop(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{
			`var sum = 0
			for i in 1..5 ::
				sum = sum + i
			end
			sum`,
			10, // 1 + 2 + 3 + 4 = 10 (exclusive, doesn't include 5)
		},
		{
			`var sum = 0
			for i in 1...5 ::
				sum = sum + i
			end
			sum`,
			15, // 1 + 2 + 3 + 4 + 5 = 15 (inclusive)
		},
		{
			`var count = 0
			for i in 0..10 ::
				count = count + 1
			end
			count`,
			10, // 0, 1, 2, 3, 4, 5, 6, 7, 8, 9 = 10 iterations
		},
		{
			`var last = 0
			for i in 5...8 ::
				last = i
			end
			last`,
			8, // Last value in inclusive range
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestRangeWithVariables(t *testing.T) {
	input := `
	var start = 1
	var finish = 5
	var sum = 0
	for i in start...finish ::
		sum = sum + i
	end
	sum
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 15) // 1 + 2 + 3 + 4 + 5
}

// String Indexing Tests

func TestStringIndexing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"[0]`, "h"},
		{`"hello"[1]`, "e"},
		{`"hello"[4]`, "o"},
		{`"world"[0]`, "w"},
		{`"a"[0]`, "a"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestStringIndexWithVariable(t *testing.T) {
	input := `
	var str = "hello"
	var i = 2
	str[i]
	`

	evaluated := testEval(input)
	testStringObject(t, evaluated, "l")
}

func TestStringIndexOutOfBounds(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`"hello"[10]`},
		{`"hello"[100]`},
		{`"hello"[-1]`},
		{`""[0]`},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNullObject(t, evaluated)
	}
}

func TestStringIndexInExpression(t *testing.T) {
	input := `
	var str = "abc"
	var first = str[0]
	var second = str[1]
	first + second
	`

	evaluated := testEval(input)
	testStringObject(t, evaluated, "ab")
}

// Break Statement Tests

func TestBreakInLoop(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{
			`var sum = 0
			for i in [1, 2, 3, 4, 5] ::
				if i == 3 ::
					break
				end
				sum = sum + i
			end
			sum`,
			3, // 1 + 2 = 3, breaks before adding 3
		},
		{
			`var count = 0
			for i in 1..100 ::
				count = count + 1
				if count == 5 ::
					break
				end
			end
			count`,
			5, // Stops after 5 iterations
		},
		{
			`var result = 0
			for i in [10, 20, 30, 40, 50] ::
				if i > 25 ::
					break
				end
				result = i
			end
			result`,
			20, // Last value before break
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestBreakInNestedLoop(t *testing.T) {
	input := `
	var count = 0
	for i in [1, 2, 3] ::
		for j in [1, 2, 3] ::
			count = count + 1
			if j == 2 ::
				break
			end
		end
	end
	count
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 6) // Each outer loop runs inner 2 times: 3 * 2 = 6
}

// String Interpolation Tests

func TestStringInterpolation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`var name = "World"
			"Hello #{name}!"`,
			"Hello World!",
		},
		{
			`var x = 5
			var y = 3
			"#{x} + #{y} = #{x + y}"`,
			"5 + 3 = 8",
		},
		{
			`var count = 42
			"The answer is #{count}"`,
			"The answer is 42",
		},
		{
			`"Result: #{10 * 2}"`,
			"Result: 20",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestStringInterpolationWithArrays(t *testing.T) {
	input := `
	var arr = [1, 2, 3]
	"First element: #{arr[0]}"
	`

	evaluated := testEval(input)
	testStringObject(t, evaluated, "First element: 1")
}

func TestStringInterpolationNested(t *testing.T) {
	input := `
	var x = 5
	var y = 10
	"x = #{x}, y = #{y}, sum = #{x + y}"
	`

	evaluated := testEval(input)
	testStringObject(t, evaluated, "x = 5, y = 10, sum = 15")
}

// Property/Dot Access Tests

func TestPropertyMethodAccess(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"hello".length`, 5},
		{`"HELLO".lower`, "hello"},
		{`"hello".upper`, "HELLO"},
		{`[1, 2, 3].length`, 3},
		{`[1, 2, 3].first`, 1},
		{`[1, 2, 3].last`, 3},
		{`42.to_string`, "42"},
		{`true.to_string`, "true"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, expected)
		case string:
			testStringObject(t, evaluated, expected)
		}
	}
}

func TestMapPropertyAccess(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`var obj = {"name": "Alice", "age": 25}
			obj.name`,
			"Alice",
		},
		{
			`var obj = {"name": "Bob", "age": 30}
			obj.age`,
			30,
		},
		{
			`var obj = {"x": 10, "y": 20}
			obj.x + obj.y`,
			30,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, expected)
		case string:
			testStringObject(t, evaluated, expected)
		}
	}
}

func TestPropertyAssignment(t *testing.T) {
	input := `
	var obj = {"count": 0}
	obj.count = 5
	obj.count
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 5)
}

func TestCustomPropertyOnPrimitives(t *testing.T) {
	input := `
	var str = "hello"
	str.custom = fn(self) :: self.upper end
	str.custom
	`

	evaluated := testEval(input)
	testStringObject(t, evaluated, "HELLO")
}

// String Comparison Tests

func TestStringComparisonOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"apple" < "banana"`, true},
		{`"banana" < "apple"`, false},
		{`"apple" > "banana"`, false},
		{`"banana" > "apple"`, true},
		{`"apple" <= "apple"`, true},
		{`"apple" <= "banana"`, true},
		{`"apple" >= "apple"`, true},
		{`"banana" >= "apple"`, true},
		{`"abc" == "abc"`, true},
		{`"abc" != "def"`, true},
		{`"hello" < "world"`, true},
		{`"zebra" > "aardvark"`, true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestStringComparisonCaseSensitive(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"Apple" == "apple"`, false},
		{`"Apple" != "apple"`, true},
		{`"A" < "a"`, true}, // Uppercase comes before lowercase in ASCII
		{`"z" > "Z"`, true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

// Array Equality Edge Cases

func TestArrayEqualityDeep(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			`var a = [1, 2, 3]
			var b = [1, 2, 3]
			check :: a is b end`,
			true,
		},
		{
			`var a = [1, 2, 3]
			var b = [1, 2, 4]
			check :: a is b end`,
			false,
		},
		{
			`var a = [[1, 2], [3, 4]]
			var b = [[1, 2], [3, 4]]
			check :: a is b end`,
			true,
		},
		{
			`var a = [[1, 2], [3, 4]]
			var b = [[1, 2], [3, 5]]
			check :: a is b end`,
			false,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testResult, ok := evaluated.(*object.TestResult)
		if !ok {
			t.Errorf("Expected TestResult, got %T", evaluated)
			continue
		}

		if tt.expected && testResult.Failed > 0 {
			t.Errorf("Expected arrays to be equal, but test failed")
		}
		if !tt.expected && testResult.Passed > 0 {
			t.Errorf("Expected arrays to not be equal, but test passed")
		}
	}
}

func TestArrayEqualityDifferentTypes(t *testing.T) {
	input := `
	var a = [1, "two", true]
	var b = [1, "two", true]
	check :: a is b end
	`

	evaluated := testEval(input)
	testResult, ok := evaluated.(*object.TestResult)
	if !ok {
		t.Fatalf("Expected TestResult, got %T", evaluated)
	}

	if testResult.Failed > 0 {
		t.Error("Expected mixed-type arrays to be equal")
	}
}

func TestArrayEqualityDifferentLengths(t *testing.T) {
	input := `
	var a = [1, 2, 3]
	var b = [1, 2]
	check :: a is b end
	`

	evaluated := testEval(input)
	testResult, ok := evaluated.(*object.TestResult)
	if !ok {
		t.Fatalf("Expected TestResult, got %T", evaluated)
	}

	if testResult.Passed > 0 {
		t.Error("Expected arrays with different lengths to not be equal")
	}
}

// Module System Tests

func TestModuleDeclaration(t *testing.T) {
	input := `
	module Math ::
		var PI = 3.14159
		fn square(x) ::
			return x * x
		end
	end
	Math
	`

	evaluated := testEval(input)
	module, ok := evaluated.(*object.Module)
	if !ok {
		t.Fatalf("object is not Module. got=%T (%+v)", evaluated, evaluated)
	}

	if module.Name != "Math" {
		t.Errorf("module name wrong. got=%q, want=%q", module.Name, "Math")
	}
}

func TestModulePropertyAccess(t *testing.T) {
	input := `
	module Math ::
		var PI = 3.14159
	end
	Math.PI
	`

	evaluated := testEval(input)
	testNumberObject(t, evaluated, 3.14159)
}

func TestModuleFunctionCall(t *testing.T) {
	input := `
	module Math ::
		fn add(a, b) ::
			return a + b
		end
	end
	Math.add(5, 3)
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 8)
}

func TestTypeAlias(t *testing.T) {
	input := `
	type MyNumber = number
	var x: MyNumber = 42
	check ::
		x isA MyNumber
	end
	`

	evaluated := testEval(input)
	testResult, ok := evaluated.(*object.TestResult)
	if !ok {
		t.Fatalf("Expected TestResult, got %T", evaluated)
	}

	if testResult.Failed > 0 {
		t.Error("Expected type alias check to pass")
	}
}

func TestNestedModules(t *testing.T) {
	input := `
	module Outer ::
		var x = 10
		module Inner ::
			var y = 20
		end
	end
	Outer.x + Outer.Inner.y
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 30)
}

// Additional Edge Cases

func TestEmptyForLoop(t *testing.T) {
	input := `
	var count = 0
	for i in [] ::
		count = count + 1
	end
	count
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 0)
}

func TestForLoopWithRangeAndBreak(t *testing.T) {
	input := `
	var sum = 0
	for i in 1...100 ::
		if i > 10 ::
			break
		end
		sum = sum + i
	end
	sum
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 55) // 1+2+3+4+5+6+7+8+9+10 = 55
}

func TestConstantReassignmentError(t *testing.T) {
	input := `
	const x = 5
	x = 10
	`

	evaluated := testEval(input)
	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("Expected error object, got %T", evaluated)
	}

	if !contains(errObj.Message, "constant") {
		t.Errorf("Expected constant reassignment error, got: %s", errObj.Message)
	}
}

func TestShadowingInNestedScopes(t *testing.T) {
	input := `
	var x = 10
	fn test() ::
		var x = 20
		return x
	end
	var inner = test()
	var outer = x
	inner + outer
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 30) // 20 (inner) + 10 (outer)
}
