package evaluator

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/vpaulo/seda/ast"
	"github.com/vpaulo/seda/lexer"
	"github.com/vpaulo/seda/object"
	"github.com/vpaulo/seda/parser"
	"github.com/vpaulo/seda/pkg"
)

// Global flag to prevent infinite recursion in where block tests
var in_where_block_test = false

// Global type objects to hold user-defined methods
var global_array_object *object.Map
var global_string_object *object.Map
var global_number_object *object.Map
var global_map_object *object.Map

func init() {
	// Initialize the global type objects
	global_array_object = &object.Map{Pairs: make(map[string]object.MapPair)}
	global_string_object = &object.Map{Pairs: make(map[string]object.MapPair)}
	global_number_object = &object.Map{Pairs: make(map[string]object.MapPair)}
	global_map_object = &object.Map{Pairs: make(map[string]object.MapPair)}

	// Set up the evaluator reference for object_methods
	SetEvaluator(func(node interface{}, env *object.Environment) object.Object {
		if ast_node, ok := node.(ast.Node); ok {
			return Eval(ast_node, env)
		}
		return object.NewError("invalid node type")
	})

	// Share the global type objects with object_methods
	SetArrayRegistry(global_array_object)
	SetStringRegistry(global_string_object)
	SetNumberRegistry(global_number_object)
	SetMapRegistry(global_map_object)
}

// Eval evaluates an AST node and returns an object
func Eval(node ast.Node, env *object.Environment) object.Object {
	// Handle nil nodes from parse errors
	if node == nil {
		return object.NewError("parse error: invalid syntax")
	}

	switch node := node.(type) {
	// Program
	case *ast.Program:
		return eval_program(node.Statements, env)

	// Statements
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if is_error(val) {
			return val
		}
		if node.IsConstant {
			env.SetConstant(node.Name.Value, val)
		} else {
			env.Set(node.Name.Value, val)
		}
		return val

	case *ast.BlockStatement:
		return eval_block_statement(node, env)

	case *ast.IfStatement:
		return eval_if_statement(node, env)

	case *ast.ForStatement:
		return eval_for_statement(node, env)

	case *ast.CaseStatement:
		return eval_case_statement(node, env)

	case *ast.FnStatement:
		return eval_fn_statement(node, env)

	case *ast.ReturnStatement:
		return eval_return_statement(node, env)

	case *ast.BreakStatement:
		return &object.Break{}

	case *ast.CheckStatement:
		return eval_check_statement(node, env)

	case *ast.ModuleStatement:
		return eval_module_statement(node, env)

	case *ast.TypeStatement:
		return eval_type_statement(node, env)

	case *ast.UsingStatement:
		return eval_using_statement(node, env)

	// Expressions
	case *ast.NumberLiteral:
		return eval_number_literal(node)

	case *ast.StringLiteral:
		return eval_string_literal(node)

	case *ast.InterpolatedString:
		return eval_interpolated_string(node, env)

	case *ast.BooleanLiteral:
		return eval_boolean_literal(node)

	case *ast.Identifier:
		return eval_identifier(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if is_error(right) {
			return right
		}
		return eval_prefix_expression(node.Operator, right)

	case *ast.InfixExpression:
		// Short-circuit evaluation for logical operators
		if node.Operator == "&&" || node.Operator == "and" {
			left := Eval(node.Left, env)
			if is_error(left) {
				return left
			}
			// Short-circuit: if left is falsy, return false without evaluating right
			if !object.IsTruthy(left) {
				return object.FALSE
			}
			// Left is truthy, evaluate right
			right := Eval(node.Right, env)
			if is_error(right) {
				return right
			}
			return native_bool(object.IsTruthy(right))
		}

		if node.Operator == "||" || node.Operator == "or" {
			left := Eval(node.Left, env)
			if is_error(left) {
				return left
			}
			// Short-circuit: if left is truthy, return true without evaluating right
			if object.IsTruthy(left) {
				return object.TRUE
			}
			// Left is falsy, evaluate right
			right := Eval(node.Right, env)
			if is_error(right) {
				return right
			}
			return native_bool(object.IsTruthy(right))
		}

		// For all other operators, evaluate both sides
		left := Eval(node.Left, env)
		if is_error(left) {
			return left
		}
		right := Eval(node.Right, env)
		if is_error(right) {
			return right
		}
		return eval_infix_expression(node.Operator, left, right)

	case *ast.ArrayLiteral:
		elements := eval_expressions(node.Elements, env)
		if len(elements) == 1 && is_error(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.MapLiteral:
		return eval_map_literal(node, env)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if is_error(left) {
			return left
		}
		index := Eval(node.Index, env)
		if is_error(index) {
			return index
		}
		return eval_index_expression(left, index)

	case *ast.AssignmentExpression:
		val := Eval(node.Value, env)
		if is_error(val) {
			return val
		}

		// Handle identifier assignment
		if ident, ok := node.Left.(*ast.Identifier); ok {
			result := env.Update(ident.Value, val)
			if is_error(result) {
				return result
			}
			return val
		}

		// Handle index assignment (e.g., map["key"] = value, array[0] = value)
		if index_expr, ok := node.Left.(*ast.IndexExpression); ok {
			// Evaluate the left side to get the collection
			collection := Eval(index_expr.Left, env)
			if is_error(collection) {
				return collection
			}

			// Evaluate the index
			index := Eval(index_expr.Index, env)
			if is_error(index) {
				return index
			}

			// Handle map index assignment
			if map_obj, ok := collection.(*object.Map); ok {
				// Convert index to string key
				var key string
				switch idx := index.(type) {
				case *object.String:
					key = idx.Value
				case *object.Number:
					key = idx.Inspect()
				default:
					key = index.Inspect()
				}

				map_obj.Pairs[key] = object.MapPair{
					Key:   index,
					Value: val,
				}
				return val
			}

			// Handle array index assignment
			if array_obj, ok := collection.(*object.Array); ok {
				if num_idx, ok := index.(*object.Number); ok {
					idx := int(num_idx.Value)
					if idx < 0 || idx >= len(array_obj.Elements) {
						return object.NewError("index out of bounds: %d", idx)
					}
					array_obj.Elements[idx] = val
					return val
				}
				return object.NewError("array index must be a number, got %s", index.Type())
			}

			return object.NewError("index assignment not supported for %s", collection.Type())
		}

		// Handle property assignment (e.g., Array.map = fn(...) :: ... end)
		if dot_expr, ok := node.Left.(*ast.DotExpression); ok {
			// Evaluate the left side to get the object
			obj := Eval(dot_expr.Left, env)
			if is_error(obj) {
				return obj
			}

			property_name := dot_expr.Property.Value

			// Check if it's a Map object
			if map_obj, ok := obj.(*object.Map); ok {
				// If the key already exists in Pairs, update it there (data update)
				// OR if the value is not a function, treat it as data
				if _, exists := map_obj.Pairs[property_name]; exists || val.Type() != object.FUNCTION_OBJ {
					map_obj.Pairs[property_name] = object.MapPair{
						Key:   &object.String{Value: property_name},
						Value: val,
					}
				} else {
					// It's a new function being added - treat as a custom method
					if map_obj.Properties == nil {
						map_obj.Properties = make(map[string]object.Object)
					}
					map_obj.Properties[property_name] = val
				}
				return val
			}

			// Check if it's an Array object
			if array_obj, ok := obj.(*object.Array); ok {
				if array_obj.Properties == nil {
					array_obj.Properties = make(map[string]object.Object)
				}
				array_obj.Properties[property_name] = val
				return val
			}

			// Check if it's a String object
			if str_obj, ok := obj.(*object.String); ok {
				if str_obj.Properties == nil {
					str_obj.Properties = make(map[string]object.Object)
				}
				str_obj.Properties[property_name] = val
				return val
			}

			// Check if it's a Number object
			if num_obj, ok := obj.(*object.Number); ok {
				if num_obj.Properties == nil {
					num_obj.Properties = make(map[string]object.Object)
				}
				num_obj.Properties[property_name] = val
				return val
			}

			// Check if it's a Boolean object
			if bool_obj, ok := obj.(*object.Boolean); ok {
				if bool_obj.Properties == nil {
					bool_obj.Properties = make(map[string]object.Object)
				}
				bool_obj.Properties[property_name] = val
				return val
			}

			return object.NewError("cannot assign property to %s", obj.Type())
		}

		return object.NewError("invalid assignment target: %T", node.Left)

	case *ast.FunctionLiteral:
		return eval_function_literal(node, env)

	case *ast.CallExpression:
		return eval_call_expression(node, env)

	case *ast.DotExpression:
		return eval_dot_expression(node, env)

	case *ast.CaseExpression:
		return eval_case_expression(node, env)

	case *ast.RangeExpression:
		return eval_range_expression(node, env)

	default:
		return object.NewError("unknown node type: %T", node)
	}
}

// eval_program evaluates a program (list of statements)
func eval_program(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

// eval_block_statement evaluates a block of statements
func eval_block_statement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ || rt == object.BREAK_OBJ {
				return result
			}
		}
	}

	return result
}

// Literal evaluation functions
func eval_number_literal(node *ast.NumberLiteral) object.Object {
	// Convert string to float64
	var value float64
	_, err := fmt.Sscanf(node.Value, "%f", &value)
	if err != nil {
		return object.NewError("invalid number: %s", node.Value)
	}
	return &object.Number{Value: value}
}

func eval_string_literal(node *ast.StringLiteral) object.Object {
	return &object.String{Value: node.Value}
}

func eval_interpolated_string(node *ast.InterpolatedString, env *object.Environment) object.Object {
	var result string

	for _, part := range node.Parts {
		// Evaluate each part
		evaluated := Eval(part, env)

		// Check for errors
		if is_error(evaluated) {
			return evaluated
		}

		// Convert to string and append
		// Use Value directly for strings to avoid quoted output
		switch val := evaluated.(type) {
		case *object.String:
			result += val.Value
		default:
			result += evaluated.Inspect()
		}
	}

	return &object.String{Value: result}
}

func eval_boolean_literal(node *ast.BooleanLiteral) object.Object {
	return native_bool(node.Value)
}

func eval_identifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		// Check for global functions (print, println)
		if builtin := get_global_function(node.Value); builtin != nil {
			return builtin
		}
		// Check for global type objects
		if node.Value == "Array" {
			if global_array_object == nil {
				global_array_object = &object.Map{Pairs: make(map[string]object.MapPair)}
			}
			return global_array_object
		}
		if node.Value == "String" {
			if global_string_object == nil {
				global_string_object = &object.Map{Pairs: make(map[string]object.MapPair)}
			}
			return global_string_object
		}
		if node.Value == "Number" {
			if global_number_object == nil {
				global_number_object = &object.Map{Pairs: make(map[string]object.MapPair)}
			}
			return global_number_object
		}
		if node.Value == "Map" {
			if global_map_object == nil {
				global_map_object = &object.Map{Pairs: make(map[string]object.MapPair)}
			}
			return global_map_object
		}
		return object.NewError("identifier not found: %s", node.Value)
	}
	return val
}

// Expression evaluation helpers
func eval_expressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if is_error(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func eval_map_literal(node *ast.MapLiteral, env *object.Environment) object.Object {
	pairs := make(map[string]object.MapPair)

	for _, pair := range node.Pairs {
		key := Eval(pair.Key, env)
		if is_error(key) {
			return key
		}

		value := Eval(pair.Value, env)
		if is_error(value) {
			return value
		}

		// Use string representation of key
		keyStr := key.String()
		pairs[keyStr] = object.MapPair{Key: key, Value: value}
	}

	return &object.Map{Pairs: pairs}
}

func eval_index_expression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.NUMBER_OBJ:
		return eval_array_index_expression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.NUMBER_OBJ:
		return eval_string_index_expression(left, index)
	case left.Type() == object.MAP_OBJ:
		return eval_map_index_expression(left, index)
	default:
		return object.NewError("index operator not supported: %s", left.Type())
	}
}

func eval_string_index_expression(str, index object.Object) object.Object {
	string_object := str.(*object.String)
	idx := int(index.(*object.Number).Value)
	max := len(string_object.Value) - 1

	if idx < 0 || idx > max {
		return object.NULL
	}

	return &object.String{Value: string(string_object.Value[idx])}
}

func eval_array_index_expression(array, index object.Object) object.Object {
	array_object := array.(*object.Array)
	idx := int(index.(*object.Number).Value)
	max := len(array_object.Elements) - 1

	if idx < 0 || idx > max {
		return object.NULL
	}

	return array_object.Elements[idx]
}

func eval_map_index_expression(map_obj, index object.Object) object.Object {
	map_object := map_obj.(*object.Map)
	key := index.String()

	pair, ok := map_object.Pairs[key]
	if !ok {
		return object.NULL
	}

	return pair.Value
}

// Prefix expression evaluation
func eval_prefix_expression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return eval_bang_operator_expression(right)
	case "-":
		return eval_minus_prefix_operator_expression(right)
	default:
		return object.NewError("unknown operator: %s%s", operator, right.Type())
	}
}

func eval_bang_operator_expression(right object.Object) object.Object {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func eval_minus_prefix_operator_expression(right object.Object) object.Object {
	if right.Type() != object.NUMBER_OBJ {
		return object.NewError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Number).Value
	return &object.Number{Value: -value}
}

// Infix expression evaluation
func eval_infix_expression(operator string, left, right object.Object) object.Object {
	// Handle logical operators with truthy conversion for any type
	if operator == "&&" || operator == "and" || operator == "||" || operator == "or" {
		return eval_logical_infix_expression(operator, left, right)
	}

	switch {
	case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
		return eval_number_infix_expression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return eval_string_infix_expression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return eval_boolean_infix_expression(operator, left, right)
	case operator == "==":
		return native_bool(left == right)
	case operator == "!=":
		return native_bool(left != right)
	case left.Type() != right.Type():
		return object.NewError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return object.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func eval_number_infix_expression(operator string, left, right object.Object) object.Object {
	left_val := left.(*object.Number).Value
	right_val := right.(*object.Number).Value

	switch operator {
	case "+":
		return &object.Number{Value: left_val + right_val}
	case "-":
		return &object.Number{Value: left_val - right_val}
	case "*":
		return &object.Number{Value: left_val * right_val}
	case "/":
		if right_val == 0 {
			return object.NewError("division by zero")
		}
		return &object.Number{Value: left_val / right_val}
	case "%":
		if right_val == 0 {
			return object.NewError("division by zero")
		}
		return &object.Number{Value: float64(int(left_val) % int(right_val))}
	case "^":
		return &object.Number{Value: math.Pow(left_val, right_val)}
	case "<":
		return native_bool(left_val < right_val)
	case ">":
		return native_bool(left_val > right_val)
	case "<=":
		return native_bool(left_val <= right_val)
	case ">=":
		return native_bool(left_val >= right_val)
	case "==":
		return native_bool(left_val == right_val)
	case "!=":
		return native_bool(left_val != right_val)
	default:
		return object.NewError("unknown operator: %s", operator)
	}
}

func eval_string_infix_expression(operator string, left, right object.Object) object.Object {
	left_val := left.(*object.String).Value
	right_val := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: left_val + right_val}
	case "==":
		return native_bool(left_val == right_val)
	case "!=":
		return native_bool(left_val != right_val)
	case "<":
		return native_bool(left_val < right_val)
	case ">":
		return native_bool(left_val > right_val)
	case "<=":
		return native_bool(left_val <= right_val)
	case ">=":
		return native_bool(left_val >= right_val)
	default:
		return object.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func eval_logical_infix_expression(operator string, left, right object.Object) object.Object {
	left_truthy := object.IsTruthy(left)
	right_truthy := object.IsTruthy(right)

	switch operator {
	case "&&", "and":
		return native_bool(left_truthy && right_truthy)
	case "||", "or":
		return native_bool(left_truthy || right_truthy)
	default:
		return object.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func eval_boolean_infix_expression(operator string, left, right object.Object) object.Object {
	left_val := left.(*object.Boolean).Value
	right_val := right.(*object.Boolean).Value

	switch operator {
	case "==":
		return native_bool(left_val == right_val)
	case "!=":
		return native_bool(left_val != right_val)
	default:
		return object.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// Utility functions
func native_bool(input bool) *object.Boolean {
	if input {
		return object.TRUE
	}
	return object.FALSE
}

func is_error(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// Control flow evaluation functions

func eval_if_statement(node *ast.IfStatement, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if is_error(condition) {
		return condition
	}

	if is_truthy(condition) {
		return Eval(node.ThenBlock, env)
	}

	// Check else-if clauses
	for _, else_if := range node.ElseIfs {
		condition := Eval(else_if.Condition, env)
		if is_error(condition) {
			return condition
		}

		if is_truthy(condition) {
			return Eval(else_if.Block, env)
		}
	}

	// Check else block
	if node.ElseBlock != nil {
		return Eval(node.ElseBlock, env)
	}

	return object.NULL
}

func eval_for_statement(node *ast.ForStatement, env *object.Environment) object.Object {
	// Evaluate the iterable expression
	iterable := Eval(node.Iterable, env)
	if is_error(iterable) {
		return iterable
	}

	// Create a new environment for the loop scope
	loop_env := object.NewEnclosedEnvironment(env)

	var result object.Object = object.NULL

	switch iter := iterable.(type) {
	case *object.Array:
		for i, element := range iter.Elements {
			// Set index variable if present
			if node.Index != nil {
				loop_env.Set(node.Index.Value, &object.Number{Value: float64(i)})
			}

			// Set value variable
			loop_env.Set(node.Variable.Value, element)

			// Execute loop body
			result = Eval(node.Body, loop_env)
			if is_error(result) {
				return result
			}

			// Handle break
			if result.Type() == object.BREAK_OBJ {
				return object.NULL
			}

			// Handle return values
			if result.Type() == object.RETURN_VALUE_OBJ {
				return result
			}
		}

	case *object.String:
		// Iterate over string characters
		for i, char := range iter.Value {
			// Set index variable if present
			if node.Index != nil {
				loop_env.Set(node.Index.Value, &object.Number{Value: float64(i)})
			}

			// Set value variable (character as string)
			loop_env.Set(node.Variable.Value, &object.String{Value: string(char)})

			// Execute loop body
			result = Eval(node.Body, loop_env)
			if is_error(result) {
				return result
			}

			// Handle break
			if result.Type() == object.BREAK_OBJ {
				return object.NULL
			}

			// Handle return values
			if result.Type() == object.RETURN_VALUE_OBJ {
				return result
			}
		}

	case *object.Range:
		// Iterate over range
		end := iter.End
		if iter.Inclusive {
			end = iter.End + 1
		}
		for i := iter.Start; i < end; i++ {
			// Set index variable if present (for ranges, this is the iteration count)
			if node.Index != nil {
				loop_env.Set(node.Index.Value, &object.Number{Value: float64(i - iter.Start)})
			}

			// Set value variable (the current number in the range)
			loop_env.Set(node.Variable.Value, &object.Number{Value: float64(i)})

			// Execute loop body
			result = Eval(node.Body, loop_env)
			if is_error(result) {
				return result
			}

			// Handle break
			if result.Type() == object.BREAK_OBJ {
				return object.NULL
			}

			// Handle return values
			if result.Type() == object.RETURN_VALUE_OBJ {
				return result
			}
		}

	case *object.Map:
		// Iterate over map key-value pairs
		for key, pair := range iter.Pairs {
			// Set index variable if present (for maps, this is the value)
			if node.Index != nil {
				loop_env.Set(node.Index.Value, pair.Value)
			}

			// Set value variable (the key)
			loop_env.Set(node.Variable.Value, &object.String{Value: key})

			// Execute loop body
			result = Eval(node.Body, loop_env)
			if is_error(result) {
				return result
			}

			// Handle break
			if result.Type() == object.BREAK_OBJ {
				return object.NULL
			}

			// Handle return values
			if result.Type() == object.RETURN_VALUE_OBJ {
				return result
			}
		}

	default:
		return object.NewError("object is not iterable: %T", iterable)
	}

	return result
}

func eval_case_statement(node *ast.CaseStatement, env *object.Environment) object.Object {
	if node == nil {
		return object.NewError("case statement is nil")
	}

	// Evaluate the expression to match against
	expr := Eval(node.Expression, env)
	if is_error(expr) {
		return expr
	}

	// Try each case branch
	for _, branch := range node.Branches {
		if branch == nil {
			continue
		}

		// Check for wildcard pattern (underscore) before evaluation
		if ident, ok := branch.Pattern.(*ast.Identifier); ok && ident.Value == "_" {
			// Wildcard matches everything
			return Eval(branch.Result, env)
		}

		// Evaluate the pattern
		pattern := Eval(branch.Pattern, env)
		if is_error(pattern) {
			return pattern
		}

		// Check if pattern matches expression
		if is_equal(expr, pattern) {
			return Eval(branch.Result, env)
		}
	}

	// No match found
	return object.NULL
}

func eval_case_expression(node *ast.CaseExpression, env *object.Environment) object.Object {
	if node == nil {
		return object.NewError("case expression is nil")
	}

	// Evaluate the expression to match against
	expr := Eval(node.Expression, env)
	if is_error(expr) {
		return expr
	}

	// Try each case branch
	for _, branch := range node.Branches {
		if branch == nil {
			continue
		}

		// Check for wildcard pattern (underscore) before evaluation
		if ident, ok := branch.Pattern.(*ast.Identifier); ok && ident.Value == "_" {
			// Wildcard matches everything
			return Eval(branch.Result, env)
		}

		// Evaluate the pattern
		pattern := Eval(branch.Pattern, env)
		if is_error(pattern) {
			return pattern
		}

		// Check if pattern matches expression
		if is_equal(expr, pattern) {
			return Eval(branch.Result, env)
		}
	}

	// No match found - return NULL for expressions
	return object.NULL
}

func eval_range_expression(node *ast.RangeExpression, env *object.Environment) object.Object {
	start := Eval(node.Start, env)
	if is_error(start) {
		return start
	}

	end := Eval(node.End, env)
	if is_error(end) {
		return end
	}

	start_num, ok := start.(*object.Number)
	if !ok {
		return object.NewError("range start must be a number, got %s", start.Type())
	}

	end_num, ok := end.(*object.Number)
	if !ok {
		return object.NewError("range end must be a number, got %s", end.Type())
	}

	return &object.Range{
		Start:     int(start_num.Value),
		End:       int(end_num.Value),
		Inclusive: node.Inclusive,
	}
}

// Helper functions for control flow

func is_truthy(obj object.Object) bool {
	switch obj {
	case object.NULL:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		return true
	}
}

func is_equal(left, right object.Object) bool {
	// Handle different types
	if left.Type() != right.Type() {
		return false
	}

	switch left.Type() {
	case object.NUMBER_OBJ:
		return left.(*object.Number).Value == right.(*object.Number).Value
	case object.STRING_OBJ:
		return left.(*object.String).Value == right.(*object.String).Value
	case object.BOOLEAN_OBJ:
		return left.(*object.Boolean).Value == right.(*object.Boolean).Value
	case object.NULL_OBJ:
		return true // both are null
	case object.ARRAY_OBJ:
		left_arr := left.(*object.Array)
		right_arr := right.(*object.Array)

		// Arrays must have same length
		if len(left_arr.Elements) != len(right_arr.Elements) {
			return false
		}

		// Compare each element recursively
		for i := range left_arr.Elements {
			if !is_equal(left_arr.Elements[i], right_arr.Elements[i]) {
				return false
			}
		}
		return true
	default:
		return left == right
	}
}

// Function evaluation functions

func eval_fn_statement(node *ast.FnStatement, env *object.Environment) object.Object {
	// Function definitions create a function object and bind it to the environment
	fn := &object.Function{
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        env,
		WhereBlock: node.WhereBlock,
	}

	// Bind the function to the environment
	env.Set(node.Name.Value, fn)
	return fn
}

func eval_return_statement(node *ast.ReturnStatement, env *object.Environment) object.Object {
	var val object.Object = object.NULL

	if node.Value != nil {
		val = Eval(node.Value, env)
		if is_error(val) {
			return val
		}
	}

	return &object.ReturnValue{Value: val}
}

func eval_function_literal(node *ast.FunctionLiteral, env *object.Environment) object.Object {
	return &object.Function{
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        env,
	}
}

func eval_call_expression(node *ast.CallExpression, env *object.Environment) object.Object {
	// Check if this is a method call (obj.method())
	if dot_expr, ok := node.Function.(*ast.DotExpression); ok {
		return eval_method_call(dot_expr, node.Arguments, env)
	}

	// Regular function call
	function := Eval(node.Function, env)
	if is_error(function) {
		return function
	}

	args := eval_expressions(node.Arguments, env)
	if len(args) == 1 && is_error(args[0]) {
		return args[0]
	}

	return apply_function(function, args, env)
}

func apply_function(fn object.Object, args []object.Object, callerEnv *object.Environment) object.Object {
	switch function := fn.(type) {
	case *object.Function:
		extended_env := extend_function_env(function, args)
		evaluated := Eval(function.Body, extended_env)
		result := unwrap_return_value(evaluated)

		// Execute where block assertions if present (but not if we're already in a where block test)
		if function.WhereBlock != nil && !in_where_block_test {
			test_result := eval_where_block(function.WhereBlock, extended_env, result, args)
			if test_result.Failed > 0 {
				// For now, print test failures but still return the function result
				fmt.Printf("Function test failures:\n%s\n", test_result.String())
			}
		}

		return result
	case *object.Builtin:
		return function.Fn(args...)
	default:
		return object.NewError("not a function: %T", fn)
	}
}

func extend_function_env(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		if paramIdx < len(args) {
			env.Set(param.Name.Value, args[paramIdx])
		}
	}

	return env
}

func unwrap_return_value(obj object.Object) object.Object {
	if return_value, ok := obj.(*object.ReturnValue); ok {
		return return_value.Value
	}
	return obj
}

func eval_dot_expression(node *ast.DotExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if is_error(left) {
		return left
	}

	// Handle module access
	if module, ok := left.(*object.Module); ok {
		if value, exists := module.Environment.Get(node.Property.Value); exists {
			return value
		}
		return object.NewError("undefined property '%s' in module '%s'", node.Property.Value, module.Name)
	}
	
	// Handle map property access - check data keys first, then custom methods
	if map_obj, ok := left.(*object.Map); ok {
		property_name := node.Property.Value
		// First check if it's a data key in Pairs
		if pair, exists := map_obj.Pairs[property_name]; exists {
			return pair.Value
		}
		// If not a data key, fall through to method call below
	}

	// Handle property-style method access (zero-argument methods without parentheses)
	// Try to call the method with no arguments
	method_name := node.Property.Value
	result := call_object_method(left, method_name, []object.Object{})

	// If it's an error saying the method doesn't exist, return that error
	// Otherwise return the result (which could be a value or an error)
	return result
}

func eval_method_call(dot_expr *ast.DotExpression, arguments []ast.Expression, env *object.Environment) object.Object {
	// Evaluate the object (receiver)
	receiver := Eval(dot_expr.Left, env)
	if is_error(receiver) {
		return receiver
	}

	// Handle module method calls
	if module, ok := receiver.(*object.Module); ok {
		// Get function from module environment
		function_name := dot_expr.Property.Value
		if function, exists := module.Environment.Get(function_name); exists {
			// Evaluate arguments
			args := eval_expressions(arguments, env)
			if len(args) == 1 && is_error(args[0]) {
				return args[0]
			}

			// Call the function
			return apply_function(function, args, env)
		}
		return object.NewError("undefined function '%s' in module '%s'", function_name, module.Name)
	}

	// Evaluate arguments
	args := eval_expressions(arguments, env)
	if len(args) == 1 && is_error(args[0]) {
		return args[0]
	}

	// Get method name
	method_name := dot_expr.Property.Value

	// Dispatch method based on receiver type
	return call_object_method(receiver, method_name, args)
}

// Testing evaluation functions

func eval_check_statement(node *ast.CheckStatement, env *object.Environment) object.Object {
	result := &object.TestResult{
		Passed:   0,
		Failed:   0,
		Failures: []string{},
		Label:    node.Label,
	}
	
	// Create a new environment for the check block so variables are scoped
	check_env := object.NewEnclosedEnvironment(env)

	// Evaluate statements (e.g., var/const declarations) first
	for _, stmt := range node.Statements {
		eval_result := Eval(stmt, check_env)
		if is_error(eval_result) {
			result.Failed++
			result.Failures = append(result.Failures, eval_result.(*object.Error).Message)
			return result
		}
	}

	// Then evaluate assertions using the check block environment
	for _, assertion := range node.Assertions {
		passed, message := eval_assertion(assertion, check_env)
		if passed {
			result.Passed++
		} else {
			result.Failed++
			result.Failures = append(result.Failures, message)
		}
	}

	return result
}

func eval_module_statement(node *ast.ModuleStatement, env *object.Environment) object.Object {
	// Create a new environment for the module
	module_env := object.NewEnclosedEnvironment(env)

	// Evaluate the module body in the module environment
	result := eval_block_statement(node.Body, module_env)
	if is_error(result) {
		return result
	}

	// Create the module object
	module := &object.Module{
		Name:        node.Name.Value,
		Environment: module_env,
	}

	// Register the module in the current environment
	env.Set(node.Name.Value, module)

	return module
}

// eval_type_statement handles type alias declarations
func eval_type_statement(node *ast.TypeStatement, env *object.Environment) object.Object {
	// Create a type alias object
	type_alias := &object.TypeAlias{
		Name:           node.Name.Value,
		TypeAnnotation: node.Type,
	}

	// Store the type alias in the environment
	env.Set(node.Name.Value, type_alias)

	return type_alias
}

func eval_using_statement(node *ast.UsingStatement, env *object.Environment) object.Object {
	// Get the module path
	module_path := node.Path.Value

	// Resolve the path (handle relative paths)
	resolved_path, err := resolve_module_path(module_path)
	if err != nil {
		return object.NewError("failed to resolve module path '%s': %s", module_path, err.Error())
	}

	// Load and parse the module file
	content, err := os.ReadFile(resolved_path)
	if err != nil {
		return object.NewError("failed to read module file '%s': %s", resolved_path, err.Error())
	}

	// Parse the module file
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if p.HasErrors() {
		errors := strings.Join(p.FormatErrors(), "; ")
		return object.NewError("parse errors in module '%s': %s", module_path, errors)
	}

	// Create a new environment for the loaded module
	module_env := object.NewEnclosedEnvironment(env)

	// Evaluate the module file in its own environment
	result := Eval(program, module_env)
	if is_error(result) {
		return result
	}

	// Extract modules from the module environment and register them in the current environment
	for name, value := range module_env.GetStore() {
		if module, ok := value.(*object.Module); ok {
			// Use alias if provided, otherwise use original module name
			module_name := name
			if node.Alias != nil {
				module_name = node.Alias.Value
			}
			env.Set(module_name, module)
		}
	}

	return &object.Null{}
}

func eval_where_block(where_block *ast.WhereBlock, env *object.Environment, return_value object.Object, args []object.Object) *object.TestResult {
	// Set global flag to prevent infinite recursion
	in_where_block_test = true
	defer func() { in_where_block_test = false }()

	result := &object.TestResult{
		Passed:   0,
		Failed:   0,
		Failures: []string{},
		Label:    "where block",
	}

	// Create a test environment with special variables
	test_env := object.NewEnclosedEnvironment(env)

	// Add special variables for testing
	test_env.Set("result", return_value)
	if len(args) > 0 {
		test_env.Set("arg0", args[0])
	}
	if len(args) > 1 {
		test_env.Set("arg1", args[1])
	}
	if len(args) > 2 {
		test_env.Set("arg2", args[2])
	}

	for _, assertion := range where_block.Assertions {
		passed, message := eval_assertion(assertion, test_env)
		if passed {
			result.Passed++
		} else {
			result.Failed++
			result.Failures = append(result.Failures, message)
		}
	}

	return result
}

func eval_assertion(assertion *ast.Assertion, env *object.Environment) (bool, string) {
	left := Eval(assertion.Left, env)
	if is_error(left) {
		return false, fmt.Sprintf("Error evaluating left side: %s", left.String())
	}

	right := Eval(assertion.Right, env)
	if is_error(right) {
		return false, fmt.Sprintf("Error evaluating right side: %s", right.String())
	}

	switch assertion.Operator {
	case "is":
		return eval_is_assertion(left, right)
	case "isA":
		return eval_isA_assertion(left, right)
	case "contains":
		return eval_contains_assertion(left, right)
	default:
		return false, fmt.Sprintf("Unknown assertion operator: %s", assertion.Operator)
	}
}

func eval_is_assertion(left, right object.Object) (bool, string) {
	equal := is_equal(left, right)
	if !equal {
		return false, fmt.Sprintf("Expected %s, got %s", right.Inspect(), left.Inspect())
	}
	return true, ""
}

func eval_isA_assertion(left, right object.Object) (bool, string) {
	var expectedType string

	// Right can be either a string type name or a TypeAlias
	switch r := right.(type) {
	case *object.String:
		expectedType = strings.ToLower(r.Value)
	case *object.TypeAlias:
		// For type aliases, check the underlying type
		expectedType = strings.ToLower(r.TypeAnnotation.Name)
	default:
		return false, "isA operator requires a string type name or type alias"
	}

	actualType := string(left.Type())

	// Convert our internal type names to user-friendly names
	userFriendlyType := get_user_friendly_type_name(actualType)

	if userFriendlyType != expectedType {
		return false, fmt.Sprintf("Expected type %s, got %s", expectedType, userFriendlyType)
	}
	return true, ""
}

func eval_contains_assertion(left, right object.Object) (bool, string) {
	switch container := left.(type) {
	case *object.Array:
		for _, element := range container.Elements {
			if is_equal(element, right) {
				return true, ""
			}
		}
		return false, fmt.Sprintf("Array %s does not contain %s", left.Inspect(), right.Inspect())

	case *object.String:
		if right.Type() != object.STRING_OBJ {
			return false, "Cannot check if string contains non-string value"
		}
		needle := right.(*object.String).Value
		haystack := container.Value
		if strings.Contains(haystack, needle) {
			return true, ""
		}
		return false, fmt.Sprintf("String %s does not contain %s", left.Inspect(), right.Inspect())

	default:
		return false, fmt.Sprintf("Cannot use contains on type %s", left.Type())
	}
}

func get_user_friendly_type_name(internal_type string) string {
	switch internal_type {
	case "NUMBER":
		return "number"
	case "STRING":
		return "string"
	case "BOOLEAN":
		return "boolean"
	case "ARRAY":
		return "array"
	case "MAP":
		return "map"
	case "FUNCTION":
		return "function"
	case "NULL":
		return "null"
	default:
		return internal_type
	}
}

// Test Runner functionality

func RunTests(program *ast.Program, env *object.Environment) *object.TestResult {
	// First, execute the program to set up all variables and functions
	Eval(program, env)

	// Then collect and run all check blocks
	total_result := &object.TestResult{
		Passed:   0,
		Failed:   0,
		Failures: []string{},
		Label:    "Test Suite",
	}

	check_blocks := collect_check_blocks(program)

	for _, check_block := range check_blocks {
		result := eval_check_statement(check_block, env)
		if test_result, ok := result.(*object.TestResult); ok {
			total_result.Passed += test_result.Passed
			total_result.Failed += test_result.Failed
			for _, failure := range test_result.Failures {
				label := test_result.Label
				if label == "" {
					label = "unnamed test"
				}
				total_result.Failures = append(total_result.Failures, fmt.Sprintf("[%s] %s", label, failure))
			}
		}
	}

	return total_result
}

func collect_check_blocks(node ast.Node) []*ast.CheckStatement {
	var check_blocks []*ast.CheckStatement

	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Statements {
			check_blocks = append(check_blocks, collect_check_blocks(stmt)...)
		}
	case *ast.CheckStatement:
		check_blocks = append(check_blocks, n)
	case *ast.BlockStatement:
		for _, stmt := range n.Statements {
			check_blocks = append(check_blocks, collect_check_blocks(stmt)...)
		}
	case *ast.IfStatement:
		check_blocks = append(check_blocks, collect_check_blocks(n.ThenBlock)...)
		for _, else_if := range n.ElseIfs {
			check_blocks = append(check_blocks, collect_check_blocks(else_if.Block)...)
		}
		if n.ElseBlock != nil {
			check_blocks = append(check_blocks, collect_check_blocks(n.ElseBlock)...)
		}
	case *ast.ForStatement:
		check_blocks = append(check_blocks, collect_check_blocks(n.Body)...)
	case *ast.FnStatement:
		check_blocks = append(check_blocks, collect_check_blocks(n.Body)...)
	}

	return check_blocks
}

// resolve_module_path resolves module paths for import with URI support
func resolve_module_path(path string) (string, error) {
	switch {
	case strings.HasPrefix(path, "std/"):
		return resolve_std_module(path)
	case strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../"):
		return resolve_relative_module(path)
	case filepath.IsAbs(path):
		return path, nil
	case strings.Contains(path, "github.com/") || strings.Contains(path, "gitlab.com/") || strings.Contains(path, "bitbucket.org/"):
		return resolve_third_party_module(path)
	case strings.Contains(path, "/"):
		// Local subdirectory path (e.g., "utils/string.s")
		return resolve_local_module(path)
	default:
		return resolve_local_module(path)
	}
}

// resolve_std_module resolves standard library modules
func resolve_std_module(path string) (string, error) {
	// Get home directory
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %v", err)
	}

	// Standard library path: ~/.seda/std/
	std_path := filepath.Join(home_dir, ".seda", "std", strings.TrimPrefix(path, "std/")) + ".s"

	// Check if file exists
	if _, err := os.Stat(std_path); os.IsNotExist(err) {
		return "", fmt.Errorf("standard library module '%s' not found at %s", path, std_path)
	}

	return std_path, nil
}

// resolve_third_party_module resolves third-party modules (e.g., github.com/user/repo)
func resolve_third_party_module(path string) (string, error) {
	// Get home directory
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %v", err)
	}

	// Third-party package path: ~/.seda/packages/github.com/user/repo/module.s
	package_path := filepath.Join(home_dir, ".seda", "packages", path, "module.s")

	// Check if already cached
	if _, err := os.Stat(package_path); err == nil {
		return package_path, nil
	}

	// Try to download if it's a git repository
	if strings.Contains(path, "github.com/") {
		return download_git_module(path, package_path)
	}

	return "", fmt.Errorf("third-party module '%s' not found and could not be downloaded", path)
}

// resolve_relative_module resolves relative paths
func resolve_relative_module(path string) (string, error) {
	abs_path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return abs_path, nil
}

// resolve_local_module resolves local modules in current directory
func resolve_local_module(path string) (string, error) {
	// Add .s extension if not present
	if !strings.HasSuffix(path, ".s") {
		path = path + ".s"
	}

	abs_path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return abs_path, nil
}

// download_git_module downloads a git repository module using the package manager
func download_git_module(path, targetPath string) (string, error) {
	// Use package manager to install the module
	manager := pkg.NewManager()

	// Convert path to repository URL
	repo_url := "https://" + path + ".git"

	// Install the package
	if err := manager.Install(repo_url); err != nil {
		return "", fmt.Errorf("could not install package %s: %v", path, err)
	}

	// Check if module.s exists in the installed package
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return "", fmt.Errorf("module.s not found in repository %s", path)
	}

	return targetPath, nil
}