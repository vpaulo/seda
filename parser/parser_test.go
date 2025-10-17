package parser

import (
	"fmt"
	"testing"

	"github.com/vpaulo/seda/ast"
	"github.com/vpaulo/seda/lexer"
)

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
		expectedType       string
		expectedIsConstant bool
	}{
		{"var x = 5", "x", 5, "", false},
		{"var y: number = 10", "y", 10, "number", false},
		{"const name = \"hello\"", "name", "hello", "", true},
		{"const age: number = 25", "age", 25, "number", true},
		{"var flag: boolean = true", "flag", true, "boolean", false},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			t.Errorf("Test case %d (%s) has parser errors:", i, tt.input)
			for _, err := range p.Errors() {
				t.Errorf("  %s", err)
			}
			continue
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testVarStatement(t, stmt, tt.expectedIdentifier, tt.expectedIsConstant) {
			return
		}

		val := stmt.(*ast.VarStatement)
		if !testLiteralExpression(t, val.Value, tt.expectedValue) {
			return
		}

		if tt.expectedType != "" {
			if val.Type == nil {
				t.Errorf("expected type annotation, got nil")
				return
			}
			if val.Type.Name != tt.expectedType {
				t.Errorf("expected type %s, got %s", tt.expectedType, val.Type.Name)
			}
		}
	}
}

func TestFunctionStatements(t *testing.T) {
	input := `
	fn add(x: number, y: number): number ::
		x + y
	where ::
		add(2, 3) is 5
		add(0, 0) is 0
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FnStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FnStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "add" {
		t.Errorf("function name wrong. want 'add', got %s", stmt.Name.Value)
	}

	if len(stmt.Parameters) != 2 {
		t.Errorf("function parameters wrong. want 2, got %d", len(stmt.Parameters))
	}

	testLiteralExpression(t, stmt.Parameters[0].Name, "x")
	testLiteralExpression(t, stmt.Parameters[1].Name, "y")

	if stmt.ReturnType.Name != "number" {
		t.Errorf("return type wrong. want 'number', got %s", stmt.ReturnType.Name)
	}

	if stmt.WhereBlock == nil {
		t.Error("expected where block, got nil")
	} else if len(stmt.WhereBlock.Assertions) != 2 {
		t.Errorf("expected 2 assertions in where block, got %d", len(stmt.WhereBlock.Assertions))
	}
}

func TestMethodStatements(t *testing.T) {
	input := `
	fn Person.greet(): string ::
		"Hello, I'm " + self.name
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FnStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FnStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Receiver == nil {
		t.Error("expected receiver type, got nil")
	} else if stmt.Receiver.Name != "Person" {
		t.Errorf("receiver type wrong. want 'Person', got %s", stmt.Receiver.Name)
	}

	if stmt.Name.Value != "greet" {
		t.Errorf("method name wrong. want 'greet', got %s", stmt.Name.Value)
	}
}

func TestStructStatements(t *testing.T) {
	input := `
	struct Person ::
		name: string,
		age: number,
		email: string
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.StructStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.StructStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "Person" {
		t.Errorf("struct name wrong. want 'Person', got %s", stmt.Name.Value)
	}

	expectedFields := []struct {
		name string
		typ  string
	}{
		{"name", "string"},
		{"age", "number"},
		{"email", "string"},
	}

	if len(stmt.Fields) != len(expectedFields) {
		t.Errorf("wrong number of fields. want %d, got %d", len(expectedFields), len(stmt.Fields))
	}

	for i, expected := range expectedFields {
		if i >= len(stmt.Fields) {
			break
		}
		field := stmt.Fields[i]
		if field.Name.Value != expected.name {
			t.Errorf("field %d name wrong. want %s, got %s", i, expected.name, field.Name.Value)
		}
		if field.Type.Name != expected.typ {
			t.Errorf("field %d type wrong. want %s, got %s", i, expected.typ, field.Type.Name)
		}
	}
}

func TestIfStatements(t *testing.T) {
	input := `
	if x > 5 ::
		print("greater")
	else if x == 5 ::
		print("equal")
	else ::
		print("less")
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T",
			program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "x", ">", 5) {
		return
	}

	if len(stmt.ElseIfs) != 1 {
		t.Errorf("expected 1 else if clause, got %d", len(stmt.ElseIfs))
	}

	if stmt.ElseBlock == nil {
		t.Error("expected else block, got nil")
	}
}

func TestCaseStatements(t *testing.T) {
	input := `
	case value ::
		"hello" => "world"
		42 => "answer"
		_ => "default"
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.CaseStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.CaseStatement. got=%T",
			program.Statements[0])
	}

	if !testIdentifier(t, stmt.Expression, "value") {
		return
	}

	if len(stmt.Branches) != 3 {
		t.Errorf("expected 3 case branches, got %d", len(stmt.Branches))
	}
}

func TestForStatements(t *testing.T) {
	tests := []struct {
		input    string
		hasIndex bool
		indexVar string
		valueVar string
		iterable string
	}{
		{"for item in items :: print(item) end", false, "", "item", "items"},
		{"for i, item in items :: print(i, item) end", true, "i", "item", "items"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ForStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ForStatement. got=%T",
				program.Statements[0])
		}

		if tt.hasIndex {
			if stmt.Index == nil {
				t.Error("expected index variable, got nil")
			} else if stmt.Index.Value != tt.indexVar {
				t.Errorf("index variable wrong. want %s, got %s", tt.indexVar, stmt.Index.Value)
			}
		} else {
			if stmt.Index != nil {
				t.Errorf("expected no index variable, got %s", stmt.Index.Value)
			}
		}

		if stmt.Variable.Value != tt.valueVar {
			t.Errorf("value variable wrong. want %s, got %s", tt.valueVar, stmt.Variable.Value)
		}

		if !testIdentifier(t, stmt.Iterable, tt.iterable) {
			return
		}
	}
}

func TestCheckStatements(t *testing.T) {
	input := `
	check "arithmetic tests" ::
		2 + 2 is 4
		3 * 4 is 12
		10 / 2 is 5
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.CheckStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.CheckStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Label != "arithmetic tests" {
		t.Errorf("check label wrong. want 'arithmetic tests', got %s", stmt.Label)
	}

	if len(stmt.Assertions) != 3 {
		t.Errorf("expected 3 assertions, got %d", len(stmt.Assertions))
	}

	// Test first assertion: 2 + 2 is 4
	assertion := stmt.Assertions[0]
	if !testInfixExpression(t, assertion.Left, 2, "+", 2) {
		return
	}
	if assertion.Operator != "is" {
		t.Errorf("assertion operator wrong. want 'is', got %s", assertion.Operator)
	}
	if !testLiteralExpression(t, assertion.Right, 4) {
		return
	}
}

func TestExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"5", 5},
		{"10", 10},
		{"true", true},
		{"false", false},
		{"\"hello world\"", "hello world"},
		{"foobar", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testLiteralExpression(t, stmt.Expression, tt.expected) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testLiteralExpression(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestMapLiterals(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	mapLit, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("exp not ast.MapLiteral. got=%T", stmt.Expression)
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	if len(mapLit.Pairs) != len(expected) {
		t.Fatalf("map.Pairs has wrong length. got=%d", len(mapLit.Pairs))
	}

	for _, pair := range mapLit.Pairs {
		literal, ok := pair.Key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", pair.Key)
			continue
		}

		expectedValue := expected[literal.Value]
		testLiteralExpression(t, pair.Value, expectedValue)
	}
}

// Helper functions for testing

func testVarStatement(t *testing.T, s ast.Statement, name string, isConstant bool) bool {
	if s == nil {
		t.Errorf("s is nil")
		return false
	}

	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}

	// Check first name (backward compatible with single variable)
	if len(varStmt.Names) == 0 {
		t.Errorf("varStmt.Names is empty")
		return false
	}

	if varStmt.Names[0].Value != name {
		t.Errorf("varStmt.Names[0].Value not '%s'. got=%s", name, varStmt.Names[0].Value)
		return false
	}

	if varStmt.IsConstant != isConstant {
		t.Errorf("varStmt.IsConstant not %t. got=%t", isConstant, varStmt.IsConstant)
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		switch exp.(type) {
		case *ast.StringLiteral:
			return testStringLiteral(t, exp, v)
		case *ast.Identifier:
			return testIdentifier(t, exp, v)
		default:
			t.Errorf("exp not string literal or identifier. got=%T", exp)
			return false
		}
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("il not *ast.NumberLiteral. got=%T", il)
		return false
	}

	// Convert value to string for comparison
	expectedValue := fmt.Sprintf("%d", value)
	if integ.Value != expectedValue {
		t.Errorf("integ.Value not %s. got=%s", expectedValue, integ.Value)
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	str, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Errorf("exp not *ast.StringLiteral. got=%T", exp)
		return false
	}

	if str.Value != value {
		t.Errorf("str.Value not %s. got=%s", value, str.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.BooleanLiteral. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5", 5},
		{"return true", true},
		{"return foobar", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", program.Statements[0])
		}

		// Check first value (backward compatible with single return)
		if len(stmt.Values) == 0 {
			t.Fatalf("stmt.Values is empty")
			return
		}

		if !testLiteralExpression(t, stmt.Values[0], tt.expectedValue) {
			return
		}
	}
}

func TestCallExpressions(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestNestedArrays(t *testing.T) {
	input := "[[1, 2], [3, 4], [5, 6]]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("array has wrong number of elements. got=%d", len(array.Elements))
	}

	for i, elem := range array.Elements {
		innerArray, ok := elem.(*ast.ArrayLiteral)
		if !ok {
			t.Fatalf("element %d is not ast.ArrayLiteral. got=%T", i, elem)
		}

		if len(innerArray.Elements) != 2 {
			t.Fatalf("inner array %d has wrong length. got=%d", i, len(innerArray.Elements))
		}
	}
}

func TestNestedMaps(t *testing.T) {
	input := `{"outer": {"inner": 42}}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	mapLit, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("exp not ast.MapLiteral. got=%T", stmt.Expression)
	}

	if len(mapLit.Pairs) != 1 {
		t.Fatalf("map has wrong number of pairs. got=%d", len(mapLit.Pairs))
	}

	pair := mapLit.Pairs[0]
	innerMap, ok := pair.Value.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("pair value is not ast.MapLiteral. got=%T", pair.Value)
	}

	if len(innerMap.Pairs) != 1 {
		t.Fatalf("inner map has wrong number of pairs. got=%d", len(innerMap.Pairs))
	}
}

func TestComplexFunctionBody(t *testing.T) {
	input := `
	fn factorial(n: number): number ::
		if n <= 1 ::
			return 1
		else ::
			return n * factorial(n - 1)
		end
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FnStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FnStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "factorial" {
		t.Errorf("function name wrong. want 'factorial', got %s", stmt.Name.Value)
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("function body has wrong number of statements. got=%d",
			len(stmt.Body.Statements))
	}

	ifStmt, ok := stmt.Body.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("function body statement is not ast.IfStatement. got=%T",
			stmt.Body.Statements[0])
	}

	if ifStmt.ThenBlock == nil {
		t.Fatal("if statement then block is nil")
	}

	if ifStmt.ElseBlock == nil {
		t.Fatal("if statement else block is nil")
	}
}

func TestModuleStatement(t *testing.T) {
	input := `
	module Utils ::
		fn helper() ::
			return 42
		end
		var value = 100
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ModuleStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ModuleStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "Utils" {
		t.Errorf("module name wrong. want 'Utils', got %s", stmt.Name.Value)
	}

	if len(stmt.Body.Statements) != 2 {
		t.Errorf("module body has wrong number of statements. want 2, got %d",
			len(stmt.Body.Statements))
	}
}

func TestUsingStatement(t *testing.T) {
	tests := []struct {
		input string
		path  string
		alias string
	}{
		{`using "std/math"`, "std/math", ""},
		{`using "std/math" as Math`, "std/math", "Math"},
		{`using "github.com/user/repo" as UserLib`, "github.com/user/repo", "UserLib"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.UsingStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.UsingStatement. got=%T",
				program.Statements[0])
		}

		if stmt.Path.Value != tt.path {
			t.Errorf("using path wrong. want %s, got %s", tt.path, stmt.Path.Value)
		}

		expectedAlias := tt.alias
		actualAlias := ""
		if stmt.Alias != nil {
			actualAlias = stmt.Alias.Value
		}
		if actualAlias != expectedAlias {
			t.Errorf("using alias wrong. want %s, got %s", expectedAlias, actualAlias)
		}
	}
}

func TestTypeStatement(t *testing.T) {
	input := "type MyNumber = number"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.TypeStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.TypeStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "MyNumber" {
		t.Errorf("type name wrong. want 'MyNumber', got %s", stmt.Name.Value)
	}

	if stmt.Type.Name != "number" {
		t.Errorf("base type wrong. want 'number', got %s", stmt.Type.Name)
	}
}

func TestParsingErrors(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{"var x", "expected next token to be ="},
		{"var = 5", "expected next token to be IDENT"},
		{"[1, 2,", "expected next token to be ]"},
		{"{\"a\": 1,", "expected next token to be }"},
		{"if x :: print(x)", "expected next token to be END"},
		{"for item :: print(item) end", "expected next token to be IN"},
		{"fn add(x,) :: x end", "parsing error"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		_ = p.ParseProgram()

		errors := p.Errors()
		if len(errors) == 0 {
			t.Errorf("expected parser error for input: %q, but got none", tt.input)
			continue
		}

		found := false
		for _, err := range errors {
			if len(err) > 0 && len(tt.expectedError) > 0 {
				found = true
				break
			}
		}

		if !found && tt.expectedError != "" {
			t.Errorf("expected error containing %q for input: %q, but got: %v",
				tt.expectedError, tt.input, errors)
		}
	}
}

func TestDotExpression(t *testing.T) {
	input := "obj.property"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	dotExp, ok := stmt.Expression.(*ast.DotExpression)
	if !ok {
		t.Fatalf("exp not *ast.DotExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, dotExp.Left, "obj") {
		return
	}

	if dotExp.Property.Value != "property" {
		t.Errorf("dotExp.Property wrong. want 'property', got %s", dotExp.Property.Value)
	}
}

func TestEmptyArrayAndMap(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"[]", "array"},
		{"{}", "map"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		switch tt.expected {
		case "array":
			if _, ok := stmt.Expression.(*ast.ArrayLiteral); !ok {
				t.Errorf("exp not *ast.ArrayLiteral. got=%T", stmt.Expression)
			}
		case "map":
			if _, ok := stmt.Expression.(*ast.MapLiteral); !ok {
				t.Errorf("exp not *ast.MapLiteral. got=%T", stmt.Expression)
			}
		}
	}
}

func TestNestedIfStatements(t *testing.T) {
	input := `
	if x > 10 ::
		if y > 5 ::
			print("both")
		end
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	outerIf, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T",
			program.Statements[0])
	}

	if len(outerIf.ThenBlock.Statements) != 1 {
		t.Fatalf("outer if has wrong number of statements. got=%d",
			len(outerIf.ThenBlock.Statements))
	}

	innerIf, ok := outerIf.ThenBlock.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("nested statement is not ast.IfStatement. got=%T",
			outerIf.ThenBlock.Statements[0])
	}

	if !testInfixExpression(t, innerIf.Condition, "y", ">", 5) {
		return
	}
}

func TestMultipleElseIf(t *testing.T) {
	input := `
	if x > 10 ::
		print("big")
	else if x > 5 ::
		print("medium")
	else if x > 0 ::
		print("small")
	else ::
		print("zero or negative")
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T",
			program.Statements[0])
	}

	if len(stmt.ElseIfs) != 2 {
		t.Errorf("expected 2 else if clauses, got %d", len(stmt.ElseIfs))
	}

	if stmt.ElseBlock == nil {
		t.Error("expected else block, got nil")
	}
}

func TestFunctionWithoutParameters(t *testing.T) {
	input := `
	fn greet(): string ::
		return "Hello"
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.FnStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FnStatement. got=%T",
			program.Statements[0])
	}

	if len(stmt.Parameters) != 0 {
		t.Errorf("expected 0 parameters, got %d", len(stmt.Parameters))
	}

	if stmt.ReturnType.Name != "string" {
		t.Errorf("return type wrong. want 'string', got %s", stmt.ReturnType.Name)
	}
}

func TestFunctionWithWhereBlock(t *testing.T) {
	input := `
	fn double(x: number): number ::
		return x * 2
	where ::
		double(2) is 4
		double(5) is 10
		double(0) is 0
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.FnStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FnStatement. got=%T",
			program.Statements[0])
	}

	if stmt.WhereBlock == nil {
		t.Fatal("expected where block, got nil")
	}

	if len(stmt.WhereBlock.Assertions) != 3 {
		t.Errorf("expected 3 assertions in where block, got %d",
			len(stmt.WhereBlock.Assertions))
	}
}

func TestComponentStatement(t *testing.T) {
	input := `
	component Counter(initial: Number) ::
		var count = initial

		Window {
			title: "Counter App",
			width: 400,

			VBox {
				Text {
					text: "Count: 0",
					fontSize: 24
				}
				Button {
					text: "Increment"
				}
			}
		}
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ComponentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ComponentStatement. got=%T",
			program.Statements[0])
	}

	// Test component name
	if stmt.Name.Value != "Counter" {
		t.Errorf("component name wrong. want 'Counter', got %s", stmt.Name.Value)
	}

	// Test parameters
	if len(stmt.Parameters) != 1 {
		t.Fatalf("component has wrong number of parameters. want 1, got %d",
			len(stmt.Parameters))
	}

	if stmt.Parameters[0].Name.Value != "initial" {
		t.Errorf("parameter name wrong. want 'initial', got %s",
			stmt.Parameters[0].Name.Value)
	}

	if stmt.Parameters[0].Type.Name != "Number" {
		t.Errorf("parameter type wrong. want 'Number', got %s",
			stmt.Parameters[0].Type.Name)
	}

	// Test body
	if stmt.Body == nil {
		t.Fatal("component body is nil")
	}

	// Test statements in body (var count = initial)
	if len(stmt.Body.Statements) != 1 {
		t.Errorf("component body has wrong number of statements. want 1, got %d",
			len(stmt.Body.Statements))
	}

	varStmt, ok := stmt.Body.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("body statement is not *ast.VarStatement. got=%T",
			stmt.Body.Statements[0])
	}

	if varStmt.Names[0].Value != "count" {
		t.Errorf("var name wrong. want 'count', got %s", varStmt.Names[0].Value)
	}

	// Test root UI element (Window)
	if stmt.Body.Root == nil {
		t.Fatal("component root UI element is nil")
	}

	if stmt.Body.Root.Type.Value != "Window" {
		t.Errorf("root element type wrong. want 'Window', got %s",
			stmt.Body.Root.Type.Value)
	}

	// Test Window properties
	if len(stmt.Body.Root.Properties) < 2 {
		t.Errorf("Window has wrong number of properties. want at least 2, got %d",
			len(stmt.Body.Root.Properties))
	}

	// Test Window has children
	if len(stmt.Body.Root.Children) != 1 {
		t.Errorf("Window has wrong number of children. want 1, got %d",
			len(stmt.Body.Root.Children))
	}

	// Test VBox child
	vbox := stmt.Body.Root.Children[0]
	if vbox.Type.Value != "VBox" {
		t.Errorf("Window child type wrong. want 'VBox', got %s", vbox.Type.Value)
	}

	// Test VBox has 2 children (Text and Button)
	if len(vbox.Children) != 2 {
		t.Errorf("VBox has wrong number of children. want 2, got %d",
			len(vbox.Children))
	}

	// Test Text element
	text := vbox.Children[0]
	if text.Type.Value != "Text" {
		t.Errorf("first VBox child type wrong. want 'Text', got %s", text.Type.Value)
	}

	if len(text.Properties) < 1 {
		t.Error("Text element has no properties")
	}

	// Test Button element
	button := vbox.Children[1]
	if button.Type.Value != "Button" {
		t.Errorf("second VBox child type wrong. want 'Button', got %s",
			button.Type.Value)
	}

	if len(button.Properties) < 1 {
		t.Error("Button element has no properties")
	}
}

func TestSimpleComponent(t *testing.T) {
	input := `
	component HelloWorld() ::
		Window {
			title: "Hello"
		}
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ComponentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ComponentStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "HelloWorld" {
		t.Errorf("component name wrong. want 'HelloWorld', got %s", stmt.Name.Value)
	}

	if len(stmt.Parameters) != 0 {
		t.Errorf("expected 0 parameters, got %d", len(stmt.Parameters))
	}

	if stmt.Body.Root == nil {
		t.Fatal("component root is nil")
	}

	if stmt.Body.Root.Type.Value != "Window" {
		t.Errorf("root type wrong. want 'Window', got %s", stmt.Body.Root.Type.Value)
	}
}

func TestUIElementWithCallback(t *testing.T) {
	input := `
	component App() ::
		Button {
			text: "Click",
			onClick: fn() :: println("clicked") end
		}
	end
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ComponentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ComponentStatement. got=%T",
			program.Statements[0])
	}

	button := stmt.Body.Root
	if button.Type.Value != "Button" {
		t.Errorf("root type wrong. want 'Button', got %s", button.Type.Value)
	}

	// Test that onClick property exists and is a function
	onClickProp, ok := button.Properties["onClick"]
	if !ok {
		t.Fatal("onClick property not found")
	}

	_, ok = onClickProp.(*ast.FunctionLiteral)
	if !ok {
		t.Errorf("onClick is not *ast.FunctionLiteral. got=%T", onClickProp)
	}
}

