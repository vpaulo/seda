package evaluator

import (
	"testing"

	"github.com/vpaulo/seda/object"
)

// JSON Module Tests

func TestJSONParseNumber(t *testing.T) {
	input := `
var data, err = JSON.parse("42")
data
`
	result := testEval(input)
	testNumberObject(t, result, 42)
}

func TestJSONParseString(t *testing.T) {
	input := `
var data, err = JSON.parse("\"hello\"")
data
`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	if str.Value != "hello" {
		t.Errorf("wrong value. expected=%s, got=%s", "hello", str.Value)
	}
}

func TestJSONParseBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`
var data, err = JSON.parse("true")
data
`, true},
		{`
var data, err = JSON.parse("false")
data
`, false},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testBooleanObject(t, result, tt.expected)
	}
}

func TestJSONParseNull(t *testing.T) {
	input := `
var data, err = JSON.parse("null")
isNull(data)
`
	result := testEval(input)
	testBooleanObject(t, result, true)
}

func TestJSONParseArray(t *testing.T) {
	input := `
var data, err = JSON.parse("[1, 2, 3]")
data
`
	result := testEval(input)

	arr, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("result is not Array. got=%T (%+v)", result, result)
	}

	expected := []float64{1, 2, 3}
	if len(arr.Elements) != len(expected) {
		t.Fatalf("array has wrong length. expected=%d, got=%d", len(expected), len(arr.Elements))
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

func TestJSONParseObject(t *testing.T) {
	input := `
var data, err = JSON.parse("{\"name\": \"Alice\", \"age\": 30}")
data.name
`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	if str.Value != "Alice" {
		t.Errorf("wrong value. expected=%s, got=%s", "Alice", str.Value)
	}
}

func TestJSONParseNestedObject(t *testing.T) {
	input := `
var data, err = JSON.parse("{\"user\": {\"name\": \"Bob\", \"age\": 25}}")
data.user.name
`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	if str.Value != "Bob" {
		t.Errorf("wrong value. expected=%s, got=%s", "Bob", str.Value)
	}
}

func TestJSONParseNestedArray(t *testing.T) {
	input := `
var data, err = JSON.parse("[[1, 2], [3, 4]]")
data[0][1]
`
	result := testEval(input)
	testNumberObject(t, result, 2)
}

func TestJSONParseInvalid(t *testing.T) {
	input := `
var data, err = JSON.parse("invalid json")
!isNull(err)
`
	result := testEval(input)
	testBooleanObject(t, result, true)
}

func TestJSONStringifyNumber(t *testing.T) {
	input := `JSON.stringify(42)`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	if str.Value != "42" {
		t.Errorf("wrong JSON. expected=%s, got=%s", "42", str.Value)
	}
}

func TestJSONStringifyString(t *testing.T) {
	input := `JSON.stringify("hello")`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	if str.Value != `"hello"` {
		t.Errorf("wrong JSON. expected=%s, got=%s", `"hello"`, str.Value)
	}
}

func TestJSONStringifyBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`JSON.stringify(true)`, "true"},
		{`JSON.stringify(false)`, "false"},
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

func TestJSONStringifyNull(t *testing.T) {
	input := `JSON.stringify(nil)`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	if str.Value != "null" {
		t.Errorf("wrong JSON. expected=%s, got=%s", "null", str.Value)
	}
}

func TestJSONStringifyArray(t *testing.T) {
	input := `JSON.stringify([1, 2, 3])`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	if str.Value != "[1,2,3]" {
		t.Errorf("wrong JSON. expected=%s, got=%s", "[1,2,3]", str.Value)
	}
}

func TestJSONStringifyObject(t *testing.T) {
	input := `JSON.stringify({"name": "Alice", "age": 30})`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	// JSON object keys may be in any order, so just check it contains the keys
	if !contains(str.Value, `"name":"Alice"`) || !contains(str.Value, `"age":30`) {
		t.Errorf("wrong JSON. got=%s", str.Value)
	}
}

func TestJSONStringifyNestedObject(t *testing.T) {
	input := `JSON.stringify({"user": {"name": "Bob", "age": 25}})`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	// Check it contains the expected keys
	if !contains(str.Value, `"user"`) || !contains(str.Value, `"name":"Bob"`) {
		t.Errorf("wrong JSON. got=%s", str.Value)
	}
}

func TestJSONStringifyWithIndent(t *testing.T) {
	input := `JSON.stringify({"name": "Alice"}, 2)`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	// Check that it's formatted with newlines (indented)
	if !contains(str.Value, "\n") {
		t.Errorf("expected indented JSON, got=%s", str.Value)
	}
}

func TestJSONRoundTrip(t *testing.T) {
	input := `
var original = {"name": "Alice", "age": 30, "active": true}
var json_str = JSON.stringify(original)
var parsed, err = JSON.parse(json_str)
parsed.name
`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	if str.Value != "Alice" {
		t.Errorf("wrong value after round trip. expected=%s, got=%s", "Alice", str.Value)
	}
}

func TestJSONComplexData(t *testing.T) {
	input := `
var data = {
  "users": [
    {"name": "Alice", "age": 30},
    {"name": "Bob", "age": 25}
  ],
  "count": 2
}
var json_str = JSON.stringify(data)
var parsed, err = JSON.parse(json_str)
parsed.users[1].name
`
	result := testEval(input)

	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("result is not String. got=%T (%+v)", result, result)
	}

	if str.Value != "Bob" {
		t.Errorf("wrong value. expected=%s, got=%s", "Bob", str.Value)
	}
}

func TestJSONErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			`JSON.parse()`,
			"wrong number of arguments for JSON.parse. got=0, want=1",
		},
		{
			`JSON.parse(123)`,
			"argument to JSON.parse must be STRING, got NUMBER",
		},
		{
			`JSON.stringify()`,
			"wrong number of arguments for JSON.stringify. got=0, want=1 or 2",
		},
		{
			`JSON.stringify(42, "invalid")`,
			"second argument to JSON.stringify must be NUMBER, got STRING",
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

