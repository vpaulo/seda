package evaluator

import (
	"math"
	"testing"

	"github.com/vpaulo/seda/object"
)

// Math Module Tests

func TestMathConstants(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"Math.PI", math.Pi},
		{"Math.E", math.E},
		{"Math.TAU", 2 * math.Pi},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestMathPow(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"Math.pow(2, 3)", 8},
		{"Math.pow(2, 10)", 1024},
		{"Math.pow(5, 2)", 25},
		{"Math.pow(10, 0)", 1},
		{"Math.pow(2, -1)", 0.5},
		{"Math.pow(4, 0.5)", 2}, // Square root
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestMathMaxMin(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"Math.max(1, 2, 3)", 3},
		{"Math.max(5, 2, 8, 1)", 8},
		{"Math.max(-1, -5, -3)", -1},
		{"Math.min(1, 2, 3)", 1},
		{"Math.min(5, 2, 8, 1)", 1},
		{"Math.min(-1, -5, -3)", -5},
		{"Math.max(42)", 42},
		{"Math.min(42)", 42},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestMathTrigonometry(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"Math.sin(0)", 0},
		{"Math.cos(0)", 1},
		{"Math.tan(0)", 0},
		{"Math.sin(Math.PI / 2)", 1},
		{"Math.cos(Math.PI)", -1},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		num := result.(*object.Number)

		// Allow small floating point differences
		diff := math.Abs(num.Value - tt.expected)
		if diff > 0.0000001 {
			t.Errorf("for %s: expected=%f, got=%f", tt.input, tt.expected, num.Value)
		}
	}
}

func TestMathInverseTrig(t *testing.T) {
	tests := []struct {
		input       string
		expectedMin float64
		expectedMax float64
	}{
		{"Math.asin(0)", -0.0001, 0.0001},
		{"Math.acos(1)", -0.0001, 0.0001},
		{"Math.atan(0)", -0.0001, 0.0001},
		{"Math.asin(1)", math.Pi/2 - 0.0001, math.Pi/2 + 0.0001},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		num := result.(*object.Number)

		if num.Value < tt.expectedMin || num.Value > tt.expectedMax {
			t.Errorf("for %s: expected between %f and %f, got %f",
				tt.input, tt.expectedMin, tt.expectedMax, num.Value)
		}
	}
}

func TestMathAtan2(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"Math.atan2(0, 1)", 0},
		{"Math.atan2(1, 0)", math.Pi / 2},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		num := result.(*object.Number)

		diff := math.Abs(num.Value - tt.expected)
		if diff > 0.0000001 {
			t.Errorf("for %s: expected=%f, got=%f", tt.input, tt.expected, num.Value)
		}
	}
}

func TestMathLogarithms(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"Math.log(Math.E)", 1},
		{"Math.log10(10)", 1},
		{"Math.log10(100)", 2},
		{"Math.log2(2)", 1},
		{"Math.log2(8)", 3},
		{"Math.exp(0)", 1},
		{"Math.exp(1)", math.E},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		num := result.(*object.Number)

		diff := math.Abs(num.Value - tt.expected)
		if diff > 0.0000001 {
			t.Errorf("for %s: expected=%f, got=%f", tt.input, tt.expected, num.Value)
		}
	}
}

func TestMathRandom(t *testing.T) {
	// Test that random() returns value in [0, 1)
	input := "Math.random()"
	result := testEval(input)
	num := result.(*object.Number)

	if num.Value < 0 || num.Value >= 1 {
		t.Errorf("Math.random() should return value in [0, 1), got %f", num.Value)
	}

	// Test multiple calls return different values
	result2 := testEval(input)
	num2 := result2.(*object.Number)

	// Very unlikely to get same value twice (but possible)
	if num.Value == num2.Value {
		t.Logf("Warning: Math.random() returned same value twice: %f", num.Value)
	}
}

func TestMathRandomInt(t *testing.T) {
	tests := []struct {
		input string
		min   int
		max   int
	}{
		{"Math.random_int(0, 10)", 0, 10},
		{"Math.random_int(1, 100)", 1, 100},
		{"Math.random_int(50, 60)", 50, 60},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		num := result.(*object.Number)

		intVal := int(num.Value)
		if intVal < tt.min || intVal >= tt.max {
			t.Errorf("for %s: expected value in [%d, %d), got %d",
				tt.input, tt.min, tt.max, intVal)
		}
	}
}

func TestMathComplexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{
			"Math.pow(2, 3) + Math.pow(3, 2)",
			17, // 8 + 9
		},
		{
			"Math.max(Math.min(10, 5), Math.min(3, 7))",
			5, // max(min(10,5), min(3,7)) = max(5, 3) = 5
		},
		{
			"var sum = Math.pow(3, 2) + Math.pow(4, 2)\nsum.sqrt()",
			5, // Pythagorean: sqrt(9 + 16) = 5
		},
		{
			"Math.log(Math.exp(5))",
			5, // log(e^5) = 5
		},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		num := result.(*object.Number)

		diff := math.Abs(num.Value - tt.expected)
		if diff > 0.0000001 {
			t.Errorf("for %s: expected=%f, got=%f", tt.input, tt.expected, num.Value)
		}
	}
}

func TestMathWithNumberMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{
			"var x = 5.7\nMath.pow(x.floor(), 2)",
			25, // floor(5.7)^2 = 5^2 = 25
		},
		{
			"var x = -3.2\nx.abs() + Math.pow(2, 3)",
			11.2, // abs(-3.2) + 8 = 3.2 + 8
		},
		{
			"Math.PI.floor()",
			3,
		},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestMathCircleArea(t *testing.T) {
	// Calculate area of circle with radius 5
	input := `
		var radius = 5
		var area = Math.PI * Math.pow(radius, 2)
		area
	`

	result := testEval(input)
	num := result.(*object.Number)

	expected := math.Pi * 25
	diff := math.Abs(num.Value - expected)
	if diff > 0.0000001 {
		t.Errorf("circle area calculation failed: expected=%f, got=%f", expected, num.Value)
	}
}

func TestMathDistanceFormula(t *testing.T) {
	// Calculate distance between two points using Math functions
	input := `
		var x1 = 0
		var y1 = 0
		var x2 = 3
		var y2 = 4

		var dx = x2 - x1
		var dy = y2 - y1
		var sum = Math.pow(dx, 2) + Math.pow(dy, 2)
		var distance = sum.sqrt()
		distance
	`

	result := testEval(input)
	testNumberObject(t, result, 5) // 3-4-5 triangle
}

func TestMathErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"Math.pow()",
			"wrong number of arguments for Math.pow. got=0, want=2",
		},
		{
			"Math.pow(2)",
			"wrong number of arguments for Math.pow. got=1, want=2",
		},
		{
			"Math.max()",
			"Math.max requires at least one argument",
		},
		{
			"Math.min()",
			"Math.min requires at least one argument",
		},
		{
			"Math.pow(\"a\", 2)",
			"first argument to Math.pow must be NUMBER, got STRING",
		},
		{
			"Math.sin(\"hello\")",
			"argument to Math.sin must be NUMBER, got STRING",
		},
		{
			"Math.random_int(10, 5)",
			"Math.random_int: min must be less than max",
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

func TestMathVariadicFunctions(t *testing.T) {
	// Test that max/min work with variable number of arguments
	tests := []struct {
		input    string
		expected float64
	}{
		{"Math.max(1, 2)", 2},
		{"Math.max(1, 2, 3)", 3},
		{"Math.max(1, 2, 3, 4, 5)", 5},
		{"Math.min(5, 4)", 4},
		{"Math.min(5, 4, 3)", 3},
		{"Math.min(5, 4, 3, 2, 1)", 1},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestMathChaining(t *testing.T) {
	// Test that Math functions can be chained with number methods
	input := `
		var result = Math.pow(2.7, 3).floor()
		result
	`

	result := testEval(input)

	// 2.7^3 = 19.683, floor = 19
	testNumberObject(t, result, 19)
}
