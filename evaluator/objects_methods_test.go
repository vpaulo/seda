package evaluator

import (
	"testing"

	"github.com/vpaulo/seda/object"
)

// String Methods Tests

func TestStringLength(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`"hello".length()`, 5},
		{`"".length()`, 0},
		{`"a".length()`, 1},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestStringCaseTransformation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello".upper()`, "HELLO"},
		{`"WORLD".lower()`, "world"},
		{`"MiXeD".upper()`, "MIXED"},
		{`"MiXeD".lower()`, "mixed"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		str, ok := result.(*object.String)
		if !ok {
			t.Fatalf("result is not String. got=%T (%+v)", result, result)
		}
		if str.Value != tt.expected {
			t.Errorf("for %s: expected=%s, got=%s", tt.input, tt.expected, str.Value)
		}
	}
}

func TestStringSubstr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello".substr(0, 5)`, "hello"},
		{`"hello".substr(1, 3)`, "ell"},
		{`"hello".substr(2, 2)`, "ll"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		str, ok := result.(*object.String)
		if !ok {
			t.Fatalf("result is not String. got=%T (%+v)", result, result)
		}
		if str.Value != tt.expected {
			t.Errorf("for %s: expected=%s, got=%s", tt.input, tt.expected, str.Value)
		}
	}
}

func TestStringSplit(t *testing.T) {
	input := `"hello,world,test".split(",")`
	result := testEval(input)

	arr, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("result is not Array. got=%T (%+v)", result, result)
	}

	if len(arr.Elements) != 3 {
		t.Fatalf("array has wrong length. got=%d", len(arr.Elements))
	}

	expected := []string{"hello", "world", "test"}
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

func TestStringTrim(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"  hello  ".trim()`, "hello"},
		{`"\t\nhello\t\n".trim()`, "hello"},
		{`"hello".trim()`, "hello"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		str, ok := result.(*object.String)
		if !ok {
			t.Fatalf("result is not String. got=%T (%+v)", result, result)
		}
		if str.Value != tt.expected {
			t.Errorf("for %s: expected=%s, got=%s", tt.input, tt.expected, str.Value)
		}
	}
}

func TestStringReplace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello world".replace("world", "universe")`, "hello universe"},
		{`"foo bar foo".replace("foo", "baz")`, "baz bar baz"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		str, ok := result.(*object.String)
		if !ok {
			t.Fatalf("result is not String. got=%T (%+v)", result, result)
		}
		if str.Value != tt.expected {
			t.Errorf("for %s: expected=%s, got=%s", tt.input, tt.expected, str.Value)
		}
	}
}

func TestStringBooleanMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"hello".starts_with("hel")`, true},
		{`"hello".starts_with("world")`, false},
		{`"hello".ends_with("llo")`, true},
		{`"hello".ends_with("xyz")`, false},
		{`"hello world".contains("o w")`, true},
		{`"hello world".contains("xyz")`, false},
		{`"".is_empty()`, true},
		{`"hello".is_empty()`, false},
		{`"   ".is_blank()`, true},
		{`"hello".is_blank()`, false},
		{`"12345".is_numeric()`, true},
		{`"123a45".is_numeric()`, false},
		{`"hello".is_alpha()`, true},
		{`"hello123".is_alpha()`, false},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testBooleanObject(t, result, tt.expected)
	}
}

func TestStringIndexOf(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`"hello world".index_of("world")`, 6},
		{`"hello world".index_of("hello")`, 0},
		{`"hello world".index_of("xyz")`, -1},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestStringReverse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello".reverse()`, "olleh"},
		{`"abc".reverse()`, "cba"},
		{`"a".reverse()`, "a"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		str, ok := result.(*object.String)
		if !ok {
			t.Fatalf("result is not String. got=%T (%+v)", result, result)
		}
		if str.Value != tt.expected {
			t.Errorf("for %s: expected=%s, got=%s", tt.input, tt.expected, str.Value)
		}
	}
}

// Array Methods Tests

func TestArrayLength(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`[1, 2, 3].length()`, 3},
		{`[].length()`, 0},
		{`["a", "b"].length()`, 2},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestArrayPushPop(t *testing.T) {
	// Test push
	input := `
		var arr = [1, 2, 3]
		var _ = arr.push(4)
		arr.length()
	`
	result := testEval(input)
	testNumberObject(t, result, 4)

	// Test pop
	input2 := `
		var arr = [1, 2, 3]
		arr.pop()
	`
	result2 := testEval(input2)
	testNumberObject(t, result2, 3)
}

func TestArrayFirstLast(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`[1, 2, 3].first()`, 1},
		{`[1, 2, 3].last()`, 3},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestArrayMap(t *testing.T) {
	input := `
		var arr = [1, 2, 3]
		var doubled = arr.map(fn(x) :: x * 2 end)
		doubled
	`
	result := testEval(input)

	arr, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("result is not Array. got=%T (%+v)", result, result)
	}

	expected := []float64{2, 4, 6}
	if len(arr.Elements) != len(expected) {
		t.Fatalf("array has wrong length. got=%d", len(arr.Elements))
	}

	for i, exp := range expected {
		num, ok := arr.Elements[i].(*object.Number)
		if !ok {
			t.Fatalf("element %d is not Number. got=%T", i, arr.Elements[i])
		}
		if num.Value != exp {
			t.Errorf("element %d: expected=%f, got=%f", i, exp, num.Value)
		}
	}
}

func TestArrayFilter(t *testing.T) {
	input := `
		var arr = [1, 2, 3, 4, 5]
		var evens = arr.filter(fn(x) :: x % 2 == 0 end)
		evens
	`
	result := testEval(input)

	arr, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("result is not Array. got=%T (%+v)", result, result)
	}

	expected := []float64{2, 4}
	if len(arr.Elements) != len(expected) {
		t.Fatalf("array has wrong length. got=%d", len(arr.Elements))
	}

	for i, exp := range expected {
		num, ok := arr.Elements[i].(*object.Number)
		if !ok {
			t.Fatalf("element %d is not Number. got=%T", i, arr.Elements[i])
		}
		if num.Value != exp {
			t.Errorf("element %d: expected=%f, got=%f", i, exp, num.Value)
		}
	}
}

func TestArrayReduce(t *testing.T) {
	input := `
		var arr = [1, 2, 3, 4]
		var sum = arr.reduce(fn(acc, x) :: acc + x end, 0)
		sum
	`
	result := testEval(input)
	testNumberObject(t, result, 10)
}

func TestArrayFind(t *testing.T) {
	input := `
		var arr = [1, 2, 3, 4, 5]
		var found = arr.find(fn(x) :: x > 3 end)
		found
	`
	result := testEval(input)
	testNumberObject(t, result, 4)
}

func TestArrayAnyAll(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`[2, 4, 6].all(fn(x) :: x % 2 == 0 end)`, true},
		{`[1, 2, 3].all(fn(x) :: x % 2 == 0 end)`, false},
		{`[1, 2, 3].any(fn(x) :: x % 2 == 0 end)`, true},
		{`[1, 3, 5].any(fn(x) :: x % 2 == 0 end)`, false},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testBooleanObject(t, result, tt.expected)
	}
}

func TestArraySort(t *testing.T) {
	input := `
		var arr = [3, 1, 4, 1, 5, 9]
		var _ = arr.sort()
		arr
	`
	result := testEval(input)

	arr, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("result is not Array. got=%T (%+v)", result, result)
	}

	expected := []float64{1, 1, 3, 4, 5, 9}
	if len(arr.Elements) != len(expected) {
		t.Fatalf("array has wrong length. got=%d", len(arr.Elements))
	}

	for i, exp := range expected {
		num, ok := arr.Elements[i].(*object.Number)
		if !ok {
			t.Fatalf("element %d is not Number. got=%T", i, arr.Elements[i])
		}
		if num.Value != exp {
			t.Errorf("element %d: expected=%f, got=%f", i, exp, num.Value)
		}
	}
}

func TestArrayReverse(t *testing.T) {
	input := `
		var arr = [1, 2, 3, 4, 5]
		var _ = arr.reverse()
		arr
	`
	result := testEval(input)

	arr, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("result is not Array. got=%T (%+v)", result, result)
	}

	expected := []float64{5, 4, 3, 2, 1}
	if len(arr.Elements) != len(expected) {
		t.Fatalf("array has wrong length. got=%d", len(arr.Elements))
	}

	for i, exp := range expected {
		num, ok := arr.Elements[i].(*object.Number)
		if !ok {
			t.Fatalf("element %d is not Number. got=%T", i, arr.Elements[i])
		}
		if num.Value != exp {
			t.Errorf("element %d: expected=%f, got=%f", i, exp, num.Value)
		}
	}
}

func TestArrayUnique(t *testing.T) {
	input := `[1, 2, 2, 3, 3, 3, 4].unique()`
	result := testEval(input)

	arr, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("result is not Array. got=%T (%+v)", result, result)
	}

	expected := []float64{1, 2, 3, 4}
	if len(arr.Elements) != len(expected) {
		t.Fatalf("array has wrong length. got=%d", len(arr.Elements))
	}

	for i, exp := range expected {
		num, ok := arr.Elements[i].(*object.Number)
		if !ok {
			t.Fatalf("element %d is not Number. got=%T", i, arr.Elements[i])
		}
		if num.Value != exp {
			t.Errorf("element %d: expected=%f, got=%f", i, exp, num.Value)
		}
	}
}

func TestArraySlice(t *testing.T) {
	input := `[1, 2, 3, 4, 5].slice(1, 4)`
	result := testEval(input)

	arr, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("result is not Array. got=%T (%+v)", result, result)
	}

	expected := []float64{2, 3, 4}
	if len(arr.Elements) != len(expected) {
		t.Fatalf("array has wrong length. got=%d", len(arr.Elements))
	}

	for i, exp := range expected {
		num, ok := arr.Elements[i].(*object.Number)
		if !ok {
			t.Fatalf("element %d is not Number. got=%T", i, arr.Elements[i])
		}
		if num.Value != exp {
			t.Errorf("element %d: expected=%f, got=%f", i, exp, num.Value)
		}
	}
}

func TestArrayFlatten(t *testing.T) {
	input := `[[1, 2], [3, 4], [5]].flatten()`
	result := testEval(input)

	arr, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("result is not Array. got=%T (%+v)", result, result)
	}

	expected := []float64{1, 2, 3, 4, 5}
	if len(arr.Elements) != len(expected) {
		t.Fatalf("array has wrong length. got=%d", len(arr.Elements))
	}

	for i, exp := range expected {
		num, ok := arr.Elements[i].(*object.Number)
		if !ok {
			t.Fatalf("element %d is not Number. got=%T", i, arr.Elements[i])
		}
		if num.Value != exp {
			t.Errorf("element %d: expected=%f, got=%f", i, exp, num.Value)
		}
	}
}

func TestArrayJoin(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`["hello", "world"].join(" ")`, "hello world"},
		{`[1, 2, 3].join(",")`, "1,2,3"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		str, ok := result.(*object.String)
		if !ok {
			t.Fatalf("result is not String. got=%T (%+v)", result, result)
		}
		if str.Value != tt.expected {
			t.Errorf("for %s: expected=%s, got=%s", tt.input, tt.expected, str.Value)
		}
	}
}

func TestArrayStatistics(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`[1, 2, 3, 4, 5].sum()`, 15},
		{`[1, 2, 3, 4, 5].average()`, 3},
		{`[1, 2, 3, 4, 5].min()`, 1},
		{`[1, 2, 3, 4, 5].max()`, 5},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestArrayContains(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`[1, 2, 3].contains(2)`, true},
		{`[1, 2, 3].contains(5)`, false},
		{`["a", "b", "c"].contains("b")`, true},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testBooleanObject(t, result, tt.expected)
	}
}

// Number Methods Tests

func TestNumberAbs(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`(-5).abs()`, 5},
		{`(5).abs()`, 5},
		{`(-3.14).abs()`, 3.14},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestNumberFloor(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`(5.7).floor()`, 5},
		{`(-5.7).floor()`, -6},
		{`(5).floor()`, 5},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestNumberCeil(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`(5.2).ceil()`, 6},
		{`(-5.2).ceil()`, -5},
		{`(5).ceil()`, 5},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestNumberRound(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`(5.5).round()`, 6},
		{`(5.4).round()`, 5},
		{`(-5.5).round()`, -6},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestNumberSqrt(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`(16).sqrt()`, 4},
		{`(25).sqrt()`, 5},
		{`(2).sqrt()`, 1.4142135623730951},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

func TestNumberMethodChaining(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`(-5.7).abs().floor()`, 5},
		{`(3.14159).floor().abs()`, 3},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNumberObject(t, result, tt.expected)
	}
}

// Error Handling Tests

func TestStringMethodErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			`"hello".substr()`,
			"wrong number of arguments for String.substr. got=0, want=2",
		},
		{
			`"hello".substr("a", 2)`,
			"first argument to String.substr must be NUMBER, got STRING",
		},
		{
			`"hello".char_at(10)`,
			"index out of bounds",
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

func TestArrayMethodErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			`[].pop()`,
			"", // Returns NULL, not error
		},
		{
			`[].average()`,
			"cannot compute average of empty array",
		},
		{
			`[1, "a", 3].sum()`,
			"Array.sum requires all elements to be NUMBER, got STRING",
		},
	}

	for _, tt := range tests {
		result := testEval(tt.input)

		if tt.expectedMessage == "" {
			// Expecting NULL
			if result != object.NULL {
				t.Errorf("Expected NULL for %s, got %T", tt.input, result)
			}
			continue
		}

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
