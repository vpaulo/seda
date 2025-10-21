package evaluator

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"github.com/vpaulo/seda/ast"
	"github.com/vpaulo/seda/lexer"
	"github.com/vpaulo/seda/object"
	"github.com/vpaulo/seda/parser"
	"github.com/vpaulo/seda/pkg"
	"github.com/vpaulo/seda/ui"
)

// Global flag to prevent infinite recursion in where block tests
var in_where_block_test = false

// Global slice to collect where block test results during test mode
var where_block_results []*object.TestResult
var in_test_mode = false

// Global type objects to hold user-defined methods
var global_array_object *object.Map
var global_string_object *object.Map
var global_number_object *object.Map
var global_map_object *object.Map

// Global module objects
var global_math_module *object.Map
var global_file_module *object.Map
var global_json_module *object.Map
var global_os_module *object.Map
var global_time_module *object.Map
var global_ui_module *object.Map

func init() {
	// Initialize the global type objects
	global_array_object = &object.Map{Pairs: make(map[string]object.MapPair)}
	global_string_object = &object.Map{Pairs: make(map[string]object.MapPair)}
	global_number_object = &object.Map{Pairs: make(map[string]object.MapPair)}
	global_map_object = &object.Map{Pairs: make(map[string]object.MapPair)}

	// Initialize global modules
	global_math_module = init_math_module()
	global_file_module = init_file_module()
	global_json_module = init_json_module()
	global_os_module = init_os_module()
	global_time_module = init_time_module()
	global_ui_module = init_ui_module()

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

	// Set up evaluator functions for UI ComponentInstance
	ui.SetEvalFunc(func(node interface{}, env *object.Environment) object.Object {
		if ast_node, ok := node.(ast.Node); ok {
			return Eval(ast_node, env)
		}
		return object.NewError("invalid node type")
	})
	ui.SetEvalUIElementFunc(func(node interface{}, env *object.Environment) object.Object {
		if ui_element, ok := node.(*ast.UIElement); ok {
			return eval_ui_element(ui_element, env)
		}
		return object.NewError("invalid UI element node")
	})
	ui.SetIsErrorFunc(is_error)
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
		// Propagate runtime errors immediately, but allow user-created errors to be assigned
		if is_runtime_error(val) {
			return val
		}

		// Single assignment (backward compatible)
		if len(node.Names) == 1 {
			if node.IsConstant {
				// Mark the value as immutable (deep immutability)
				mark_immutable(val)
				env.SetConstant(node.Names[0].Value, val)
			} else {
				env.Set(node.Names[0].Value, val)
			}
			return val
		}

		// Multiple assignment (destructuring)
		var values []object.Object

		// Check if value is MultiValue
		if mv, ok := val.(*object.MultiValue); ok {
			values = mv.Values
		} else {
			// Single value assigned to multiple variables - error
			return object.NewError("cannot assign single value to %d variables", len(node.Names))
		}

		// Check counts match
		if len(values) != len(node.Names) {
			return object.NewError("assignment count mismatch: %d values for %d variables",
				len(values), len(node.Names))
		}

		// Assign each value to corresponding variable
		for i, name := range node.Names {
			if node.IsConstant {
				// Mark each value as immutable (deep immutability)
				mark_immutable(values[i])
				env.SetConstant(name.Value, values[i])
			} else {
				env.Set(name.Value, values[i])
			}
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

	case *ast.ComponentStatement:
		return eval_component_statement(node, env)

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

	case *ast.NilLiteral:
		return object.NULL

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

	case *ast.UIElement:
		return eval_ui_element(node, env)

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
				// Check if map is immutable
				if map_obj.IsImmutable {
					return object.NewError("cannot modify immutable map")
				}

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
				// Check if array is immutable
				if array_obj.IsImmutable {
					return object.NewError("cannot modify immutable array")
				}

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
			// Only propagate runtime errors immediately
			// User-created errors (via error() builtin) are treated as regular values
			if !result.IsUserCreated {
				return result
			}
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
			if rt == object.RETURN_VALUE_OBJ || rt == object.BREAK_OBJ {
				return result
			}
			// Propagate runtime errors immediately, but not user-created errors
			if rt == object.ERROR_OBJ {
				if err, ok := result.(*object.Error); ok && !err.IsUserCreated {
					return result
				}
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

		// Propagate runtime errors immediately
		// User-created errors can be interpolated into strings
		if is_runtime_error(evaluated) {
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
		// Check for global modules
		if node.Value == "Math" {
			if global_math_module == nil {
				global_math_module = init_math_module()
			}
			return global_math_module
		}
		if node.Value == "File" {
			if global_file_module == nil {
				global_file_module = init_file_module()
			}
			return global_file_module
		}
		if node.Value == "JSON" {
			if global_json_module == nil {
				global_json_module = init_json_module()
			}
			return global_json_module
		}
		if node.Value == "OS" {
			if global_os_module == nil {
				global_os_module = init_os_module()
			}
			return global_os_module
		}
		if node.Value == "Time" {
			if global_time_module == nil {
				global_time_module = init_time_module()
			}
			return global_time_module
		}
		if node.Value == "UI" {
			if global_ui_module == nil {
				global_ui_module = init_ui_module()
			}
			return global_ui_module
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
		// Propagate runtime errors immediately, but allow user-created errors as function arguments
		if is_runtime_error(evaluated) {
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

func eval_ui_element(node *ast.UIElement, env *object.Environment) object.Object {
	// Create runtime UI element
	element := &object.UIElement{
		ElementType: node.Type.Value,
		Properties:  make(map[string]object.Object),
		Children:    []*object.UIElement{},
	}

	// Evaluate properties
	for key, expr := range node.Properties {
		value := Eval(expr, env)
		if is_error(value) {
			return value
		}
		element.Properties[key] = value
	}

	// Recursively evaluate children
	for _, child := range node.Children {
		childObj := eval_ui_element(child, env)
		if is_error(childObj) {
			return childObj
		}
		// Type assertion to convert object.Object to *object.UIElement
		if childElement, ok := childObj.(*object.UIElement); ok {
			element.Children = append(element.Children, childElement)
		} else {
			return object.NewError("child is not a UI element: %s", childObj.Type())
		}
	}

	return element
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

// mark_immutable marks an object and all nested structures as immutable (deep immutability)
func mark_immutable(obj object.Object) {
	switch v := obj.(type) {
	case *object.Array:
		v.IsImmutable = true
		// Recursively mark all nested elements as immutable
		for _, elem := range v.Elements {
			mark_immutable(elem)
		}
	case *object.Map:
		v.IsImmutable = true
		// Recursively mark all nested values as immutable
		for _, pair := range v.Pairs {
			mark_immutable(pair.Value)
		}
	case *object.Number:
		v.IsImmutable = true
	case *object.String:
		v.IsImmutable = true
	case *object.Boolean:
		v.IsImmutable = true
	// Other types (NULL, functions, etc.) don't need immutability tracking
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

// is_runtime_error checks if an error should propagate immediately (runtime errors)
// Returns false for user-created errors (via error() builtin) which should be treated as values
func is_runtime_error(obj object.Object) bool {
	if err, ok := obj.(*object.Error); ok {
		return !err.IsUserCreated
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
		// Use epsilon-based comparison for floating-point numbers
		left_val := left.(*object.Number).Value
		right_val := right.(*object.Number).Value
		epsilon := 1e-8 // Tolerance for floating-point comparison (10 nanounits)
		return math.Abs(left_val-right_val) < epsilon
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

func eval_component_statement(node *ast.ComponentStatement, env *object.Environment) object.Object {
	// Component definitions create a UIComponent object and bind it to the environment
	component := &object.UIComponent{
		Name:       node.Name.Value,
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        env, // Capture closure environment
	}

	// Bind the component to the environment
	env.Set(node.Name.Value, component)
	return component
}

func eval_return_statement(node *ast.ReturnStatement, env *object.Environment) object.Object {
	// No return values - return null
	if len(node.Values) == 0 {
		return &object.ReturnValue{Value: object.NULL}
	}

	// Single return value (backward compatible)
	if len(node.Values) == 1 {
		val := Eval(node.Values[0], env)
		// Propagate runtime errors immediately
		if is_runtime_error(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}

	// Multiple return values
	// For Go-style error handling, allow user-created Error objects to be returned as values
	// But propagate runtime errors immediately
	values := []object.Object{}
	for _, expr := range node.Values {
		val := Eval(expr, env)
		if is_runtime_error(val) {
			return val
		}
		values = append(values, val)
	}

	return &object.ReturnValue{
		Value: &object.MultiValue{Values: values},
	}
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
	// Propagate runtime errors immediately, but allow user-created errors as arguments
	if len(args) == 1 && is_runtime_error(args[0]) {
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
				// Print test failures during normal execution
				if !in_test_mode {
					fmt.Printf("Function test failures:\n%s\n", test_result.String())
				}
			}
			// Collect where block results during test mode
			if in_test_mode {
				where_block_results = append(where_block_results, test_result)
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
	// Only propagate runtime errors, allow user-created errors to have methods called on them
	if is_runtime_error(left) {
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
	// Only propagate runtime errors, allow user-created errors to have methods called on them
	if is_runtime_error(receiver) {
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

	// Handle Map function calls - check Pairs for functions (like Math module functions)
	if map_obj, ok := receiver.(*object.Map); ok {
		method_name := dot_expr.Property.Value
		if pair, exists := map_obj.Pairs[method_name]; exists {
			// Evaluate arguments
			args := eval_expressions(arguments, env)
			if len(args) == 1 && is_error(args[0]) {
				return args[0]
			}

			// Call the function
			return apply_function(pair.Value, args, env)
		}
		// If not found in Pairs, fall through to call_object_method for custom methods
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
		Passed:     0,
		Failed:     0,
		Failures:   []string{},
		Label:      node.Label,
		Assertions: []object.AssertionResult{},
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

		// Create assertion result with source code
		assertionResult := object.AssertionResult{
			Passed:  passed,
			Message: message,
			Source:  assertion.String(), // Use AST string representation as source
		}
		result.Assertions = append(result.Assertions, assertionResult)

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
	resolved_path, err := resolve_module_path(module_path, env.SourceDir)
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
		Passed:     0,
		Failed:     0,
		Failures:   []string{},
		Label:      "where block",
		Assertions: []object.AssertionResult{},
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

	// Evaluate statements (e.g., var/const declarations) first
	for _, stmt := range where_block.Statements {
		eval_result := Eval(stmt, test_env)
		if is_error(eval_result) {
			result.Failed++
			result.Failures = append(result.Failures, eval_result.(*object.Error).Message)
			return result
		}
	}

	// Then evaluate assertions
	for _, assertion := range where_block.Assertions {
		passed, message := eval_assertion(assertion, test_env)

		// Create assertion result with source code
		assertionResult := object.AssertionResult{
			Passed:  passed,
			Message: message,
			Source:  assertion.String(), // Use AST string representation as source
		}
		result.Assertions = append(result.Assertions, assertionResult)

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
	// Special handling for 'raises' assertion - we want to capture errors, not propagate them
	if assertion.Operator == "raises" {
		// Evaluate left side - if it's an error, that's what we're testing for
		left := Eval(assertion.Left, env)
		// For raises, we DON'T treat errors as failures - they're expected

		// Evaluate right side (optional error message matcher)
		var right object.Object
		if assertion.Right != nil {
			right = Eval(assertion.Right, env)
			if is_error(right) {
				return false, fmt.Sprintf("Error evaluating right side: %s", right.String())
			}
		}

		return eval_raises_assertion(left, right)
	}

	// For all other assertions, evaluate left side normally
	left := Eval(assertion.Left, env)
	if is_error(left) {
		return false, fmt.Sprintf("Error evaluating left side: %s", left.String())
	}

	// For unary assertions, right will be nil
	var right object.Object
	if assertion.Right != nil {
		right = Eval(assertion.Right, env)
		if is_error(right) {
			return false, fmt.Sprintf("Error evaluating right side: %s", right.String())
		}
	}

	switch assertion.Operator {
	case "is":
		return eval_is_assertion(left, right)
	case "isA":
		return eval_isA_assertion(left, right)
	case "isNot":
		return eval_isNot_assertion(left, right)
	case "contains":
		return eval_contains_assertion(left, right)
	case "isGreater":
		return eval_isGreater_assertion(left, right)
	case "isLess":
		return eval_isLess_assertion(left, right)
	case "isTrue":
		return eval_isTrue_assertion(left)
	case "isFalse":
		return eval_isFalse_assertion(left)
	case "isEmpty":
		return eval_isEmpty_assertion(left)
	case "startsWith":
		return eval_startsWith_assertion(left, right)
	case "endsWith":
		return eval_endsWith_assertion(left, right)
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

// New assertion evaluation functions

func eval_isNot_assertion(left, right object.Object) (bool, string) {
	equal := is_equal(left, right)
	if equal {
		return false, fmt.Sprintf("Expected %s to not equal %s", left.Inspect(), right.Inspect())
	}
	return true, ""
}

func eval_isGreater_assertion(left, right object.Object) (bool, string) {
	if left.Type() != object.NUMBER_OBJ || right.Type() != object.NUMBER_OBJ {
		return false, "isGreater requires both operands to be numbers"
	}
	left_val := left.(*object.Number).Value
	right_val := right.(*object.Number).Value
	if left_val > right_val {
		return true, ""
	}
	return false, fmt.Sprintf("Expected %s to be greater than %s", left.Inspect(), right.Inspect())
}

func eval_isLess_assertion(left, right object.Object) (bool, string) {
	if left.Type() != object.NUMBER_OBJ || right.Type() != object.NUMBER_OBJ {
		return false, "isLess requires both operands to be numbers"
	}
	left_val := left.(*object.Number).Value
	right_val := right.(*object.Number).Value
	if left_val < right_val {
		return true, ""
	}
	return false, fmt.Sprintf("Expected %s to be less than %s", left.Inspect(), right.Inspect())
}

func eval_isTrue_assertion(left object.Object) (bool, string) {
	if left.Type() != object.BOOLEAN_OBJ {
		return false, fmt.Sprintf("isTrue requires a boolean, got %s", left.Type())
	}
	if left.(*object.Boolean).Value {
		return true, ""
	}
	return false, "Expected true, got false"
}

func eval_isFalse_assertion(left object.Object) (bool, string) {
	if left.Type() != object.BOOLEAN_OBJ {
		return false, fmt.Sprintf("isFalse requires a boolean, got %s", left.Type())
	}
	if !left.(*object.Boolean).Value {
		return true, ""
	}
	return false, "Expected false, got true"
}

func eval_isEmpty_assertion(left object.Object) (bool, string) {
	switch obj := left.(type) {
	case *object.Array:
		if len(obj.Elements) == 0 {
			return true, ""
		}
		return false, fmt.Sprintf("Expected array to be empty, but has %d elements", len(obj.Elements))
	case *object.String:
		if len(obj.Value) == 0 {
			return true, ""
		}
		return false, fmt.Sprintf("Expected string to be empty, but has length %d", len(obj.Value))
	default:
		return false, fmt.Sprintf("isEmpty requires an array or string, got %s", left.Type())
	}
}

func eval_startsWith_assertion(left, right object.Object) (bool, string) {
	if left.Type() != object.STRING_OBJ || right.Type() != object.STRING_OBJ {
		return false, "startsWith requires both operands to be strings"
	}
	left_str := left.(*object.String).Value
	right_str := right.(*object.String).Value
	if strings.HasPrefix(left_str, right_str) {
		return true, ""
	}
	return false, fmt.Sprintf("Expected %s to start with %s", left.Inspect(), right.Inspect())
}

func eval_endsWith_assertion(left, right object.Object) (bool, string) {
	if left.Type() != object.STRING_OBJ || right.Type() != object.STRING_OBJ {
		return false, "endsWith requires both operands to be strings"
	}
	left_str := left.(*object.String).Value
	right_str := right.(*object.String).Value
	if strings.HasSuffix(left_str, right_str) {
		return true, ""
	}
	return false, fmt.Sprintf("Expected %s to end with %s", left.Inspect(), right.Inspect())
}

func eval_raises_assertion(left, right object.Object) (bool, string) {
	// Check if left is an error
	if left.Type() != object.ERROR_OBJ {
		return false, fmt.Sprintf("Expected an error to be raised, but got %s", left.Type())
	}

	// If right is provided, check if error message matches
	if right != nil && right.Type() == object.STRING_OBJ {
		error_msg := left.(*object.Error).Message
		expected_msg := right.(*object.String).Value
		if strings.Contains(error_msg, expected_msg) {
			return true, ""
		}
		return false, fmt.Sprintf("Error message %q does not contain %q", error_msg, expected_msg)
	}

	// If no specific message required, just check that an error was raised
	return true, ""
}

// Test Runner functionality

func RunTests(program *ast.Program, env *object.Environment) *object.TestResult {
	// Enable test mode and reset where block results
	in_test_mode = true
	where_block_results = []*object.TestResult{}
	defer func() { in_test_mode = false }()

	// First, execute the program to set up all variables and functions
	// This will also execute where blocks and collect their results
	Eval(program, env)

	// Then collect and run all check blocks
	total_result := &object.TestResult{
		Passed:     0,
		Failed:     0,
		Failures:   []string{},
		Label:      "Test Suite",
		Assertions: []object.AssertionResult{},
	}

	check_blocks := collect_check_blocks(program)

	for _, check_block := range check_blocks {
		result := eval_check_statement(check_block, env)
		if test_result, ok := result.(*object.TestResult); ok {
			total_result.Passed += test_result.Passed
			total_result.Failed += test_result.Failed

			// Aggregate assertions from each test
			total_result.Assertions = append(total_result.Assertions, test_result.Assertions...)

			for _, failure := range test_result.Failures {
				label := test_result.Label
				if label == "" {
					label = "unnamed test"
				}
				total_result.Failures = append(total_result.Failures, fmt.Sprintf("[%s] %s", label, failure))
			}
		}
	}

	// Aggregate where block results
	for _, where_result := range where_block_results {
		total_result.Passed += where_result.Passed
		total_result.Failed += where_result.Failed

		// Aggregate assertions from where blocks
		total_result.Assertions = append(total_result.Assertions, where_result.Assertions...)

		for _, failure := range where_result.Failures {
			label := where_result.Label
			if label == "" {
				label = "unnamed where block"
			}
			total_result.Failures = append(total_result.Failures, fmt.Sprintf("[%s] %s", label, failure))
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
func resolve_module_path(path string, sourceDir string) (string, error) {
	switch {
	case strings.HasPrefix(path, "std/"):
		return resolve_std_module(path)
	case strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../"):
		return resolve_relative_module(path, sourceDir)
	case filepath.IsAbs(path):
		return path, nil
	case strings.Contains(path, "github.com/") || strings.Contains(path, "gitlab.com/") || strings.Contains(path, "bitbucket.org/"):
		return resolve_third_party_module(path)
	case strings.Contains(path, "/"):
		// Local subdirectory path (e.g., "utils/string.s")
		return resolve_local_module(path, sourceDir)
	default:
		return resolve_local_module(path, sourceDir)
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

	// Parse the path to extract repository and subdirectory
	// Format: github.com/user/repo/subdir/...
	repo_path, subdir := parse_git_url(path)

	// Construct the local package path
	var package_path string
	if subdir != "" {
		package_path = filepath.Join(home_dir, ".seda", "packages", repo_path, subdir, "module.s")
	} else {
		package_path = filepath.Join(home_dir, ".seda", "packages", repo_path, "module.s")
	}

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
func resolve_relative_module(path string, sourceDir string) (string, error) {
	// Resolve relative to the source file's directory
	var abs_path string
	if sourceDir != "" {
		abs_path = filepath.Join(sourceDir, path)
	} else {
		// Fallback to current working directory
		var err error
		abs_path, err = filepath.Abs(path)
		if err != nil {
			return "", err
		}
	}
	return abs_path, nil
}

// resolve_local_module resolves local modules relative to source file directory
func resolve_local_module(path string, sourceDir string) (string, error) {
	// Add .s extension if not present
	if !strings.HasSuffix(path, ".s") {
		path = path + ".s"
	}

	// Resolve relative to the source file's directory
	var abs_path string
	if sourceDir != "" {
		abs_path = filepath.Join(sourceDir, path)
	} else {
		// Fallback to current working directory
		var err error
		abs_path, err = filepath.Abs(path)
		if err != nil {
			return "", err
		}
	}
	return abs_path, nil
}

// parse_git_url parses a git URL path to extract repository and subdirectory
// Input: "github.com/user/repo/subdir/more" -> Output: ("github.com/user/repo", "subdir/more")
// Input: "github.com/user/repo" -> Output: ("github.com/user/repo", "")
func parse_git_url(path string) (repo string, subdir string) {
	// Remove trailing slashes
	path = strings.TrimSuffix(path, "/")

	// Split the path by "/"
	parts := strings.Split(path, "/")

	// For GitHub URLs, format is: github.com/user/repo/[subdir/...]
	// We need at least 3 parts: github.com, user, repo
	if len(parts) < 3 {
		return path, ""
	}

	// Repository is the first 3 parts (github.com/user/repo)
	repo = strings.Join(parts[:3], "/")

	// Subdirectory is everything after
	if len(parts) > 3 {
		subdir = strings.Join(parts[3:], "/")
	}

	return repo, subdir
}

// download_git_module downloads a git repository module using the package manager
func download_git_module(path, target_path string) (string, error) {
	// Use package manager to install the module
	manager := pkg.NewManager()

	// Parse the path to get repository and subdirectory
	repo_path, _ := parse_git_url(path)

	// Convert repository path to URL
	repo_url := "https://" + repo_path + ".git"

	// Install the package
	if err := manager.Install(repo_url); err != nil {
		return "", fmt.Errorf("could not install package %s: %v", repo_path, err)
	}

	// Check if module.s exists in the installed package (including subdirectory)
	if _, err := os.Stat(target_path); os.IsNotExist(err) {
		return "", fmt.Errorf("module.s not found in repository %s at path %s", path, target_path)
	}

	return target_path, nil
}

// init_math_module initializes the Math module with all math functions and constants
func init_math_module() *object.Map {
	math_module := &object.Map{Pairs: make(map[string]object.MapPair)}

	// Math.pow(base, exponent) - power function
	math_module.Pairs["pow"] = object.MapPair{
		Key: &object.String{Value: "pow"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return object.NewError("wrong number of arguments for Math.pow. got=%d, want=2", len(args))
				}
				base, ok := args[0].(*object.Number)
				if !ok {
					return object.NewError("first argument to Math.pow must be NUMBER, got %s", args[0].Type())
				}
				exp, ok := args[1].(*object.Number)
				if !ok {
					return object.NewError("second argument to Math.pow must be NUMBER, got %s", args[1].Type())
				}
				return &object.Number{Value: math.Pow(base.Value, exp.Value)}
			},
		},
	}

	// Math.max(...values) - variadic maximum
	math_module.Pairs["max"] = object.MapPair{
		Key: &object.String{Value: "max"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) == 0 {
					return object.NewError("Math.max requires at least one argument")
				}
				max_val := math.Inf(-1)
				for i, arg := range args {
					num, ok := arg.(*object.Number)
					if !ok {
						return object.NewError("argument %d to Math.max must be NUMBER, got %s", i, arg.Type())
					}
					if num.Value > max_val {
						max_val = num.Value
					}
				}
				return &object.Number{Value: max_val}
			},
		},
	}

	// Math.min(...values) - variadic minimum
	math_module.Pairs["min"] = object.MapPair{
		Key: &object.String{Value: "min"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) == 0 {
					return object.NewError("Math.min requires at least one argument")
				}
				min_val := math.Inf(1)
				for i, arg := range args {
					num, ok := arg.(*object.Number)
					if !ok {
						return object.NewError("argument %d to Math.min must be NUMBER, got %s", i, arg.Type())
					}
					if num.Value < min_val {
						min_val = num.Value
					}
				}
				return &object.Number{Value: min_val}
			},
		},
	}

	// Trigonometry functions
	math_module.Pairs["sin"] = object.MapPair{
		Key:   &object.String{Value: "sin"},
		Value: &object.Builtin{Fn: math_unary_builtin(math.Sin, "Math.sin")},
	}

	math_module.Pairs["cos"] = object.MapPair{
		Key:   &object.String{Value: "cos"},
		Value: &object.Builtin{Fn: math_unary_builtin(math.Cos, "Math.cos")},
	}

	math_module.Pairs["tan"] = object.MapPair{
		Key:   &object.String{Value: "tan"},
		Value: &object.Builtin{Fn: math_unary_builtin(math.Tan, "Math.tan")},
	}

	math_module.Pairs["asin"] = object.MapPair{
		Key:   &object.String{Value: "asin"},
		Value: &object.Builtin{Fn: math_unary_builtin(math.Asin, "Math.asin")},
	}

	math_module.Pairs["acos"] = object.MapPair{
		Key:   &object.String{Value: "acos"},
		Value: &object.Builtin{Fn: math_unary_builtin(math.Acos, "Math.acos")},
	}

	math_module.Pairs["atan"] = object.MapPair{
		Key:   &object.String{Value: "atan"},
		Value: &object.Builtin{Fn: math_unary_builtin(math.Atan, "Math.atan")},
	}

	math_module.Pairs["atan2"] = object.MapPair{
		Key: &object.String{Value: "atan2"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return object.NewError("wrong number of arguments for Math.atan2. got=%d, want=2", len(args))
				}
				y, ok := args[0].(*object.Number)
				if !ok {
					return object.NewError("first argument to Math.atan2 must be NUMBER, got %s", args[0].Type())
				}
				x, ok := args[1].(*object.Number)
				if !ok {
					return object.NewError("second argument to Math.atan2 must be NUMBER, got %s", args[1].Type())
				}
				return &object.Number{Value: math.Atan2(y.Value, x.Value)}
			},
		},
	}

	// Logarithm functions
	math_module.Pairs["log"] = object.MapPair{
		Key:   &object.String{Value: "log"},
		Value: &object.Builtin{Fn: math_unary_builtin(math.Log, "Math.log")},
	}

	math_module.Pairs["log10"] = object.MapPair{
		Key:   &object.String{Value: "log10"},
		Value: &object.Builtin{Fn: math_unary_builtin(math.Log10, "Math.log10")},
	}

	math_module.Pairs["log2"] = object.MapPair{
		Key:   &object.String{Value: "log2"},
		Value: &object.Builtin{Fn: math_unary_builtin(math.Log2, "Math.log2")},
	}

	math_module.Pairs["exp"] = object.MapPair{
		Key:   &object.String{Value: "exp"},
		Value: &object.Builtin{Fn: math_unary_builtin(math.Exp, "Math.exp")},
	}

	// Random functions
	math_module.Pairs["random"] = object.MapPair{
		Key: &object.String{Value: "random"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewError("wrong number of arguments for Math.random. got=%d, want=0", len(args))
				}
				return &object.Number{Value: rand.Float64()}
			},
		},
	}

	math_module.Pairs["random_int"] = object.MapPair{
		Key: &object.String{Value: "random_int"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return object.NewError("wrong number of arguments for Math.random_int. got=%d, want=2", len(args))
				}
				min, ok := args[0].(*object.Number)
				if !ok {
					return object.NewError("first argument to Math.random_int must be NUMBER, got %s", args[0].Type())
				}
				max, ok := args[1].(*object.Number)
				if !ok {
					return object.NewError("second argument to Math.random_int must be NUMBER, got %s", args[1].Type())
				}
				min_int := int(min.Value)
				max_int := int(max.Value)
				if min_int >= max_int {
					return object.NewError("Math.random_int: min must be less than max")
				}
				// Generate random integer in [min, max)
				random_val := min_int + rand.Intn(max_int-min_int)
				return &object.Number{Value: float64(random_val)}
			},
		},
	}

	// Math constants
	math_module.Pairs["PI"] = object.MapPair{
		Key:   &object.String{Value: "PI"},
		Value: &object.Number{Value: math.Pi},
	}

	math_module.Pairs["E"] = object.MapPair{
		Key:   &object.String{Value: "E"},
		Value: &object.Number{Value: math.E},
	}

	math_module.Pairs["TAU"] = object.MapPair{
		Key:   &object.String{Value: "TAU"},
		Value: &object.Number{Value: 2 * math.Pi},
	}

	return math_module
}

// math_unary_builtin creates a builtin wrapper for unary math functions
func math_unary_builtin(fn func(float64) float64, name string) func(...object.Object) object.Object {
	return func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for %s. got=%d, want=1", name, len(args))
		}
		num, ok := args[0].(*object.Number)
		if !ok {
			return object.NewError("argument to %s must be NUMBER, got %s", name, args[0].Type())
		}
		return &object.Number{Value: fn(num.Value)}
	}
}

// init_file_module initializes the File module with all file and directory operations
func init_file_module() *object.Map {
	file_module := &object.Map{Pairs: make(map[string]object.MapPair)}

	// File.read(path) - read file contents, returns (content, error)
	file_module.Pairs["read"] = object.MapPair{
		Key: &object.String{Value: "read"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.read. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.read must be STRING, got %s", args[0].Type())
				}

				content, err := os.ReadFile(path.Value)
				if err != nil {
					// Return (nil, error)
					user_error := &object.Error{Message: err.Error(), IsUserCreated: true}
					return &object.MultiValue{Values: []object.Object{object.NULL, user_error}}
				}

				// Return (content, nil)
				return &object.MultiValue{Values: []object.Object{
					&object.String{Value: string(content)},
					object.NULL,
				}}
			},
		},
	}

	// File.read_lines(path) - read file as array of lines, returns (lines, error)
	file_module.Pairs["read_lines"] = object.MapPair{
		Key: &object.String{Value: "read_lines"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.read_lines. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.read_lines must be STRING, got %s", args[0].Type())
				}

				content, err := os.ReadFile(path.Value)
				if err != nil {
					// Return (nil, error)
					user_error := &object.Error{Message: err.Error(), IsUserCreated: true}
					return &object.MultiValue{Values: []object.Object{object.NULL, user_error}}
				}

				// Split content into lines
				lines := strings.Split(string(content), "\n")
				elements := make([]object.Object, len(lines))
				for i, line := range lines {
					elements[i] = &object.String{Value: line}
				}

				// Return (lines_array, nil)
				return &object.MultiValue{Values: []object.Object{
					&object.Array{Elements: elements},
					object.NULL,
				}}
			},
		},
	}

	// File.write(path, content) - write content to file, returns error or nil
	file_module.Pairs["write"] = object.MapPair{
		Key: &object.String{Value: "write"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return object.NewError("wrong number of arguments for File.write. got=%d, want=2", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("first argument to File.write must be STRING, got %s", args[0].Type())
				}
				content, ok := args[1].(*object.String)
				if !ok {
					return object.NewError("second argument to File.write must be STRING, got %s", args[1].Type())
				}

				err := os.WriteFile(path.Value, []byte(content.Value), 0644)
				if err != nil {
					return &object.Error{Message: err.Error(), IsUserCreated: true}
				}

				return object.NULL
			},
		},
	}

	// File.append(path, content) - append content to file, returns error or nil
	file_module.Pairs["append"] = object.MapPair{
		Key: &object.String{Value: "append"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return object.NewError("wrong number of arguments for File.append. got=%d, want=2", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("first argument to File.append must be STRING, got %s", args[0].Type())
				}
				content, ok := args[1].(*object.String)
				if !ok {
					return object.NewError("second argument to File.append must be STRING, got %s", args[1].Type())
				}

				f, err := os.OpenFile(path.Value, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return &object.Error{Message: err.Error(), IsUserCreated: true}
				}
				defer f.Close()

				if _, err := f.WriteString(content.Value); err != nil {
					return &object.Error{Message: err.Error(), IsUserCreated: true}
				}

				return object.NULL
			},
		},
	}

	// File.delete(path) - delete file, returns error or nil
	file_module.Pairs["delete"] = object.MapPair{
		Key: &object.String{Value: "delete"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.delete. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.delete must be STRING, got %s", args[0].Type())
				}

				err := os.Remove(path.Value)
				if err != nil {
					return &object.Error{Message: err.Error(), IsUserCreated: true}
				}

				return object.NULL
			},
		},
	}

	// File.exists(path) - check if file or directory exists
	file_module.Pairs["exists"] = object.MapPair{
		Key: &object.String{Value: "exists"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.exists. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.exists must be STRING, got %s", args[0].Type())
				}

				_, err := os.Stat(path.Value)
				if err == nil {
					return object.TRUE
				}
				if os.IsNotExist(err) {
					return object.FALSE
				}
				// Other error occurred
				return object.FALSE
			},
		},
	}

	// File.size(path) - get file size in bytes, returns (size, error)
	file_module.Pairs["size"] = object.MapPair{
		Key: &object.String{Value: "size"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.size. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.size must be STRING, got %s", args[0].Type())
				}

				info, err := os.Stat(path.Value)
				if err != nil {
					user_error := &object.Error{Message: err.Error(), IsUserCreated: true}
					return &object.MultiValue{Values: []object.Object{object.NULL, user_error}}
				}

				return &object.MultiValue{Values: []object.Object{
					&object.Number{Value: float64(info.Size())},
					object.NULL,
				}}
			},
		},
	}

	// File.is_file(path) - check if path is a file
	file_module.Pairs["is_file"] = object.MapPair{
		Key: &object.String{Value: "is_file"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.is_file. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.is_file must be STRING, got %s", args[0].Type())
				}

				info, err := os.Stat(path.Value)
				if err != nil {
					return object.FALSE
				}

				if info.IsDir() {
					return object.FALSE
				}
				return object.TRUE
			},
		},
	}

	// File.is_dir(path) - check if path is a directory
	file_module.Pairs["is_dir"] = object.MapPair{
		Key: &object.String{Value: "is_dir"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.is_dir. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.is_dir must be STRING, got %s", args[0].Type())
				}

				info, err := os.Stat(path.Value)
				if err != nil {
					return object.FALSE
				}

				if info.IsDir() {
					return object.TRUE
				}
				return object.FALSE
			},
		},
	}

	// File.list_dir(path) - list directory contents, returns (files_array, error)
	file_module.Pairs["list_dir"] = object.MapPair{
		Key: &object.String{Value: "list_dir"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.list_dir. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.list_dir must be STRING, got %s", args[0].Type())
				}

				entries, err := os.ReadDir(path.Value)
				if err != nil {
					user_error := &object.Error{Message: err.Error(), IsUserCreated: true}
					return &object.MultiValue{Values: []object.Object{object.NULL, user_error}}
				}

				elements := make([]object.Object, len(entries))
				for i, entry := range entries {
					elements[i] = &object.String{Value: entry.Name()}
				}

				return &object.MultiValue{Values: []object.Object{
					&object.Array{Elements: elements},
					object.NULL,
				}}
			},
		},
	}

	// File.mkdir(path) - create directory, returns error or nil
	file_module.Pairs["mkdir"] = object.MapPair{
		Key: &object.String{Value: "mkdir"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.mkdir. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.mkdir must be STRING, got %s", args[0].Type())
				}

				err := os.Mkdir(path.Value, 0755)
				if err != nil {
					return &object.Error{Message: err.Error(), IsUserCreated: true}
				}

				return object.NULL
			},
		},
	}

	// File.mkdir_all(path) - create directory and all parent directories, returns error or nil
	file_module.Pairs["mkdir_all"] = object.MapPair{
		Key: &object.String{Value: "mkdir_all"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.mkdir_all. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.mkdir_all must be STRING, got %s", args[0].Type())
				}

				err := os.MkdirAll(path.Value, 0755)
				if err != nil {
					return &object.Error{Message: err.Error(), IsUserCreated: true}
				}

				return object.NULL
			},
		},
	}

	// File.remove_dir(path) - remove directory, returns error or nil
	file_module.Pairs["remove_dir"] = object.MapPair{
		Key: &object.String{Value: "remove_dir"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.remove_dir. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.remove_dir must be STRING, got %s", args[0].Type())
				}

				err := os.RemoveAll(path.Value)
				if err != nil {
					return &object.Error{Message: err.Error(), IsUserCreated: true}
				}

				return object.NULL
			},
		},
	}

	// File.join(...paths) - join path segments
	file_module.Pairs["join"] = object.MapPair{
		Key: &object.String{Value: "join"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) == 0 {
					return object.NewError("File.join requires at least one argument")
				}

				paths := make([]string, len(args))
				for i, arg := range args {
					str, ok := arg.(*object.String)
					if !ok {
						return object.NewError("argument %d to File.join must be STRING, got %s", i, arg.Type())
					}
					paths[i] = str.Value
				}

				result := filepath.Join(paths...)
				return &object.String{Value: result}
			},
		},
	}

	// File.basename(path) - get base name of path
	file_module.Pairs["basename"] = object.MapPair{
		Key: &object.String{Value: "basename"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.basename. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.basename must be STRING, got %s", args[0].Type())
				}

				result := filepath.Base(path.Value)
				return &object.String{Value: result}
			},
		},
	}

	// File.dirname(path) - get directory name of path
	file_module.Pairs["dirname"] = object.MapPair{
		Key: &object.String{Value: "dirname"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.dirname. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.dirname must be STRING, got %s", args[0].Type())
				}

				result := filepath.Dir(path.Value)
				return &object.String{Value: result}
			},
		},
	}

	// File.extname(path) - get file extension
	file_module.Pairs["extname"] = object.MapPair{
		Key: &object.String{Value: "extname"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.extname. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.extname must be STRING, got %s", args[0].Type())
				}

				result := filepath.Ext(path.Value)
				return &object.String{Value: result}
			},
		},
	}

	// File.absolute_path(path) - get absolute path
	file_module.Pairs["absolute_path"] = object.MapPair{
		Key: &object.String{Value: "absolute_path"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.absolute_path. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.absolute_path must be STRING, got %s", args[0].Type())
				}

				result, err := filepath.Abs(path.Value)
				if err != nil {
					return object.NewError("failed to get absolute path: %s", err.Error())
				}

				return &object.String{Value: result}
			},
		},
	}

	// File.cwd() - get current working directory
	file_module.Pairs["cwd"] = object.MapPair{
		Key: &object.String{Value: "cwd"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewError("wrong number of arguments for File.cwd. got=%d, want=0", len(args))
				}

				result, err := os.Getwd()
				if err != nil {
					return object.NewError("failed to get current directory: %s", err.Error())
				}

				return &object.String{Value: result}
			},
		},
	}

	// File.chdir(path) - change current working directory
	file_module.Pairs["chdir"] = object.MapPair{
		Key: &object.String{Value: "chdir"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for File.chdir. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to File.chdir must be STRING, got %s", args[0].Type())
				}

				err := os.Chdir(path.Value)
				if err != nil {
					return &object.Error{Message: err.Error(), IsUserCreated: true}
				}

				return object.NULL
			},
		},
	}

	return file_module
}

// init_json_module creates and returns the JSON module with parse and stringify functions
func init_json_module() *object.Map {
	json_module := &object.Map{Pairs: make(map[string]object.MapPair)}

	// JSON.parse(json_string) - parse JSON string to object
	json_module.Pairs["parse"] = object.MapPair{
		Key: &object.String{Value: "parse"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for JSON.parse. got=%d, want=1", len(args))
				}
				json_str, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to JSON.parse must be STRING, got %s", args[0].Type())
				}

				// Parse JSON into interface{}
				var data interface{}
				err := json.Unmarshal([]byte(json_str.Value), &data)
				if err != nil {
					user_error := &object.Error{Message: fmt.Sprintf("invalid JSON: %s", err.Error()), IsUserCreated: true}
					return &object.MultiValue{Values: []object.Object{object.NULL, user_error}}
				}

				// Convert Go interface{} to Seda object
				result := convert_json_to_object(data)
				return &object.MultiValue{Values: []object.Object{result, object.NULL}}
			},
		},
	}

	// JSON.stringify(obj, indent?) - convert object to JSON string
	json_module.Pairs["stringify"] = object.MapPair{
		Key: &object.String{Value: "stringify"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) < 1 || len(args) > 2 {
					return object.NewError("wrong number of arguments for JSON.stringify. got=%d, want=1 or 2", len(args))
				}

				// Convert Seda object to Go interface{}
				data := convert_object_to_json(args[0])

				var json_bytes []byte
				var err error

				// Check for optional indent parameter
				if len(args) == 2 {
					indent_obj, ok := args[1].(*object.Number)
					if !ok {
						return object.NewError("second argument to JSON.stringify must be NUMBER, got %s", args[1].Type())
					}
					indent_size := int(indent_obj.Value)
					indent_str := strings.Repeat(" ", indent_size)
					json_bytes, err = json.MarshalIndent(data, "", indent_str)
				} else {
					json_bytes, err = json.Marshal(data)
				}

				if err != nil {
					return object.NewError("failed to stringify object: %s", err.Error())
				}

				return &object.String{Value: string(json_bytes)}
			},
		},
	}

	return json_module
}

// convert_json_to_object converts a Go interface{} (from json.Unmarshal) to a Seda object
func convert_json_to_object(data interface{}) object.Object {
	switch v := data.(type) {
	case nil:
		return object.NULL
	case bool:
		if v {
			return object.TRUE
		}
		return object.FALSE
	case float64:
		return &object.Number{Value: v}
	case string:
		return &object.String{Value: v}
	case []interface{}:
		// Convert to Seda array
		elements := make([]object.Object, len(v))
		for i, elem := range v {
			elements[i] = convert_json_to_object(elem)
		}
		return &object.Array{Elements: elements}
	case map[string]interface{}:
		// Convert to Seda Map
		pairs := make(map[string]object.MapPair)
		for key, value := range v {
			pairs[key] = object.MapPair{
				Key:   &object.String{Value: key},
				Value: convert_json_to_object(value),
			}
		}
		return &object.Map{Pairs: pairs}
	default:
		return object.NewError("unsupported JSON type: %T", v)
	}
}

// convert_object_to_json converts a Seda object to a Go interface{} (for json.Marshal)
func convert_object_to_json(obj object.Object) interface{} {
	switch v := obj.(type) {
	case *object.Null:
		return nil
	case *object.Boolean:
		return v.Value
	case *object.Number:
		return v.Value
	case *object.String:
		return v.Value
	case *object.Array:
		// Convert Seda array to Go slice
		result := make([]interface{}, len(v.Elements))
		for i, elem := range v.Elements {
			result[i] = convert_object_to_json(elem)
		}
		return result
	case *object.Map:
		// Convert Seda Map to Go map
		result := make(map[string]interface{})
		for key, pair := range v.Pairs {
			result[key] = convert_object_to_json(pair.Value)
		}
		return result
	default:
		// For unsupported types, convert to string representation
		return obj.Inspect()
	}
}

// Global variable to store command line arguments
var command_line_args []string

// SetCommandLineArgs sets the command line arguments for OS.args()
func SetCommandLineArgs(args []string) {
	command_line_args = args
}

// init_os_module creates and returns the OS module with environment, process, and system functions
func init_os_module() *object.Map {
	os_module := &object.Map{Pairs: make(map[string]object.MapPair)}

	// OS.getenv(name) - get environment variable
	os_module.Pairs["getenv"] = object.MapPair{
		Key: &object.String{Value: "getenv"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for OS.getenv. got=%d, want=1", len(args))
				}
				name, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to OS.getenv must be STRING, got %s", args[0].Type())
				}

				value := os.Getenv(name.Value)
				return &object.String{Value: value}
			},
		},
	}

	// OS.setenv(name, value) - set environment variable
	os_module.Pairs["setenv"] = object.MapPair{
		Key: &object.String{Value: "setenv"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return object.NewError("wrong number of arguments for OS.setenv. got=%d, want=2", len(args))
				}
				name, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("first argument to OS.setenv must be STRING, got %s", args[0].Type())
				}
				value, ok := args[1].(*object.String)
				if !ok {
					return object.NewError("second argument to OS.setenv must be STRING, got %s", args[1].Type())
				}

				err := os.Setenv(name.Value, value.Value)
				if err != nil {
					return &object.Error{Message: err.Error(), IsUserCreated: true}
				}

				return object.NULL
			},
		},
	}

	// OS.env() - get all environment variables as map
	os_module.Pairs["env"] = object.MapPair{
		Key: &object.String{Value: "env"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewError("wrong number of arguments for OS.env. got=%d, want=0", len(args))
				}

				env_map := &object.Map{Pairs: make(map[string]object.MapPair)}
				for _, env_var := range os.Environ() {
					parts := strings.SplitN(env_var, "=", 2)
					if len(parts) == 2 {
						key := parts[0]
						value := parts[1]
						env_map.Pairs[key] = object.MapPair{
							Key:   &object.String{Value: key},
							Value: &object.String{Value: value},
						}
					}
				}

				return env_map
			},
		},
	}

	// OS.args() - get command line arguments
	os_module.Pairs["args"] = object.MapPair{
		Key: &object.String{Value: "args"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewError("wrong number of arguments for OS.args. got=%d, want=0", len(args))
				}

				elements := make([]object.Object, len(command_line_args))
				for i, arg := range command_line_args {
					elements[i] = &object.String{Value: arg}
				}

				return &object.Array{Elements: elements}
			},
		},
	}

	// OS.exit(code) - exit with status code
	os_module.Pairs["exit"] = object.MapPair{
		Key: &object.String{Value: "exit"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for OS.exit. got=%d, want=1", len(args))
				}
				code, ok := args[0].(*object.Number)
				if !ok {
					return object.NewError("argument to OS.exit must be NUMBER, got %s", args[0].Type())
				}

				os.Exit(int(code.Value))
				return object.NULL // Never reached
			},
		},
	}

	// OS.pid() - get process ID
	os_module.Pairs["pid"] = object.MapPair{
		Key: &object.String{Value: "pid"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewError("wrong number of arguments for OS.pid. got=%d, want=0", len(args))
				}

				return &object.Number{Value: float64(os.Getpid())}
			},
		},
	}

	// OS.platform() - get operating system
	os_module.Pairs["platform"] = object.MapPair{
		Key: &object.String{Value: "platform"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewError("wrong number of arguments for OS.platform. got=%d, want=0", len(args))
				}

				return &object.String{Value: runtime.GOOS}
			},
		},
	}

	// OS.arch() - get architecture
	os_module.Pairs["arch"] = object.MapPair{
		Key: &object.String{Value: "arch"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewError("wrong number of arguments for OS.arch. got=%d, want=0", len(args))
				}

				return &object.String{Value: runtime.GOARCH}
			},
		},
	}

	// OS.hostname() - get machine hostname
	os_module.Pairs["hostname"] = object.MapPair{
		Key: &object.String{Value: "hostname"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewError("wrong number of arguments for OS.hostname. got=%d, want=0", len(args))
				}

				hostname, err := os.Hostname()
				if err != nil {
					return &object.Error{Message: err.Error(), IsUserCreated: true}
				}

				return &object.String{Value: hostname}
			},
		},
	}

	// OS.home_dir() - get user home directory
	os_module.Pairs["home_dir"] = object.MapPair{
		Key: &object.String{Value: "home_dir"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewError("wrong number of arguments for OS.home_dir. got=%d, want=0", len(args))
				}

				home, err := os.UserHomeDir()
				if err != nil {
					return &object.Error{Message: err.Error(), IsUserCreated: true}
				}

				return &object.String{Value: home}
			},
		},
	}

	// OS.temp_dir() - get temporary directory
	os_module.Pairs["temp_dir"] = object.MapPair{
		Key: &object.String{Value: "temp_dir"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewError("wrong number of arguments for OS.temp_dir. got=%d, want=0", len(args))
				}

				return &object.String{Value: os.TempDir()}
			},
		},
	}

	// OS.cwd() - get current working directory
	os_module.Pairs["cwd"] = object.MapPair{
		Key: &object.String{Value: "cwd"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewError("wrong number of arguments for OS.cwd. got=%d, want=0", len(args))
				}

				cwd, err := os.Getwd()
				if err != nil {
					return &object.Error{Message: err.Error(), IsUserCreated: true}
				}

				return &object.String{Value: cwd}
			},
		},
	}

	// OS.chdir(path) - change current working directory
	os_module.Pairs["chdir"] = object.MapPair{
		Key: &object.String{Value: "chdir"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments for OS.chdir. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("argument to OS.chdir must be STRING, got %s", args[0].Type())
				}

				err := os.Chdir(path.Value)
				if err != nil {
					return &object.Error{Message: err.Error(), IsUserCreated: true}
				}

				return object.NULL
			},
		},
	}

	// OS.exec(command, ...args) - execute command and return output
	os_module.Pairs["exec"] = object.MapPair{
		Key: &object.String{Value: "exec"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) < 1 {
					return object.NewError("wrong number of arguments for OS.exec. got=%d, want=1 or more", len(args))
				}
				command, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("first argument to OS.exec must be STRING, got %s", args[0].Type())
				}

				// Convert remaining arguments to strings
				cmd_args := make([]string, len(args)-1)
				for i := 1; i < len(args); i++ {
					arg_str, ok := args[i].(*object.String)
					if !ok {
						return object.NewError("argument %d to OS.exec must be STRING, got %s", i+1, args[i].Type())
					}
					cmd_args[i-1] = arg_str.Value
				}

				// Execute command
				cmd := exec.Command(command.Value, cmd_args...)
				output, err := cmd.CombinedOutput()

				if err != nil {
					user_error := &object.Error{Message: fmt.Sprintf("command failed: %s", err.Error()), IsUserCreated: true}
					return &object.MultiValue{Values: []object.Object{object.NULL, user_error}}
				}

				return &object.MultiValue{Values: []object.Object{
					&object.String{Value: string(output)},
					object.NULL,
				}}
			},
		},
	}

	// OS.spawn(command, ...args) - spawn background process and return PID
	os_module.Pairs["spawn"] = object.MapPair{
		Key: &object.String{Value: "spawn"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) < 1 {
					return object.NewError("wrong number of arguments for OS.spawn. got=%d, want=1 or more", len(args))
				}
				command, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("first argument to OS.spawn must be STRING, got %s", args[0].Type())
				}

				// Convert remaining arguments to strings
				cmd_args := make([]string, len(args)-1)
				for i := 1; i < len(args); i++ {
					arg_str, ok := args[i].(*object.String)
					if !ok {
						return object.NewError("argument %d to OS.spawn must be STRING, got %s", i+1, args[i].Type())
					}
					cmd_args[i-1] = arg_str.Value
				}

				// Spawn command
				cmd := exec.Command(command.Value, cmd_args...)
				err := cmd.Start()

				if err != nil {
					return &object.Error{Message: fmt.Sprintf("failed to spawn command: %s", err.Error()), IsUserCreated: true}
				}

				return &object.Number{Value: float64(cmd.Process.Pid)}
			},
		},
	}

	return os_module
}

// init_time_module creates and returns the Time module
func init_time_module() *object.Map {
	time_module := &object.Map{
		Pairs: make(map[string]object.MapPair),
	}

	// Time.now() - returns current time
	time_module.Pairs["now"] = object.MapPair{
		Key: &object.String{Value: "now"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewError("Time.now() takes no arguments, got %d", len(args))
				}
				return &object.Time{Value: time.Now()}
			},
		},
	}

	// Time.date(format, dateString) - parses a date string using the given format
	// Format uses Seda patterns: YYYY, MM, DD, HH, mm, ss
	// Example: Time.date("DD-MM-YYYY", "20-01-2000")
	time_module.Pairs["date"] = object.MapPair{
		Key: &object.String{Value: "date"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return object.NewError("Time.date() takes 2 arguments (format, dateString), got %d", len(args))
				}

				// Get format string
				formatStr, ok := args[0].(*object.String)
				if !ok {
					return object.NewError("Time.date() first argument must be STRING, got %s", args[0].Type())
				}

				// Get date string
				dateStr, ok := args[1].(*object.String)
				if !ok {
					return object.NewError("Time.date() second argument must be STRING, got %s", args[1].Type())
				}

				// Convert Seda format to Go format
				goFormat := convertSedaFormatToGo(formatStr.Value)

				// Parse the date string
				parsedTime, err := time.Parse(goFormat, dateStr.Value)
				if err != nil {
					return object.NewError("Time.date() failed to parse date: %s", err.Error())
				}

				return &object.Time{Value: parsedTime}
			},
		},
	}

	// Time.unix(seconds) - creates a Time from Unix timestamp (seconds since epoch)
	time_module.Pairs["unix"] = object.MapPair{
		Key: &object.String{Value: "unix"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("Time.unix() takes 1 argument (seconds), got %d", len(args))
				}

				// Get seconds
				seconds, ok := args[0].(*object.Number)
				if !ok {
					return object.NewError("Time.unix() argument must be NUMBER, got %s", args[0].Type())
				}

				// Create time from Unix timestamp
				unixTime := time.Unix(int64(seconds.Value), 0)
				return &object.Time{Value: unixTime}
			},
		},
	}

	return time_module
}

// UI Module - Declarative UI utilities
func init_ui_module() *object.Map {
	ui_module := &object.Map{
		Pairs: make(map[string]object.MapPair),
	}

	// UI.mount(component, ...args) - instantiates and renders a component
	ui_module.Pairs["mount"] = object.MapPair{
		Key: &object.String{Value: "mount"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) == 0 {
					return object.NewError("UI.mount() requires at least one argument (component)")
				}

				// Check if first argument is a UIComponent
				component, ok := args[0].(*object.UIComponent)
				if !ok {
					return object.NewError("UI.mount() first argument must be a component, got %s", args[0].Type())
				}

				// Create the Fyne application
				app := ui.NewApplication()

				// Create event handler that can invoke Seda functions and trigger re-renders
				var componentInstance *ui.ComponentInstance
				eventHandler := func(callback object.Object) error {
					// Check if callback is a function
					fn, ok := callback.(*object.Function)
					if !ok {
						return fmt.Errorf("onClick callback is not a function")
					}

					// Invoke the function with its captured environment (closure)
					// The function should have captured the component environment
					result := apply_function(fn, []object.Object{}, fn.Env)
					if is_error(result) {
						return fmt.Errorf("error in onClick handler: %s", result.(*object.Error).Message)
					}

					// Trigger re-render after the event handler completes
					if componentInstance != nil {
						componentInstance.Rerender()
					}

					return nil
				}

				// Create renderer with app reference
				renderer := ui.NewRenderer(app, eventHandler)

				// Instantiate the component with app and renderer
				// Note: We need to set componentInstance BEFORE ShowAndRun (which blocks)
				_, err := instantiateComponent(component, args[1:], app, renderer, &componentInstance)
				if err != nil {
					return err
				}

				return object.NULL
			},
		},
	}

	// UI.inspect(ui_element) - inspects a UI element tree (debugging utility)
	ui_module.Pairs["inspect"] = object.MapPair{
		Key: &object.String{Value: "inspect"},
		Value: &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("UI.inspect() takes 1 argument, got %d", len(args))
				}

				// Check if argument is a UIElement
				element, ok := args[0].(*object.UIElement)
				if !ok {
					return object.NewError("UI.inspect() argument must be a UI element, got %s", args[0].Type())
				}

				// Return the string representation of the UI element
				return &object.String{Value: element.Inspect()}
			},
		},
	}

	return ui_module
}

// instantiateComponent creates a ComponentInstance that can re-render
func instantiateComponent(component *object.UIComponent, args []object.Object, app fyne.App, renderer *ui.Renderer, instancePtr **ui.ComponentInstance) (*ui.ComponentInstance, object.Object) {
	// Create component instance
	instance := ui.NewComponentInstance(component, args)
	instance.SetApp(app)
	instance.SetRenderer(renderer)

	// Bind parameters to arguments
	if len(args) != len(component.Parameters) {
		return nil, object.NewError("component '%s' expects %d arguments, got %d",
			component.Name, len(component.Parameters), len(args))
	}

	for i, param := range component.Parameters {
		instance.Env.Set(param.Name.Value, args[i])
	}

	// Extract window title from component body root (before full evaluation)
	// We need to check the AST directly for the title property
	windowTitle := "Seda Application" // default
	if component.Body.Root != nil {
		uiExpr := component.Body.Root
		if uiExpr.Type.Value == "Window" {
			if titleExpr, ok := uiExpr.Properties["title"]; ok {
				// Evaluate just the title property
				titleObj := Eval(titleExpr, instance.Env)
				if titleStr, ok := titleObj.(*object.String); ok {
					windowTitle = titleStr.Value
				}
			}
		}
	}

	// Create the main window
	window := app.NewWindow(windowTitle)
	window.SetPadded(true) // Add padding for better appearance
	instance.SetWindow(window)

	// Perform initial render
	err := instance.RenderComponent()
	if err != nil {
		return nil, object.NewError("failed to render component: %s", err.Error())
	}

	// Set the instance pointer BEFORE ShowAndRun (which blocks)
	// This allows event handlers to access the component instance
	*instancePtr = instance

	// Show and run the window (blocking call)
	window.ShowAndRun()

	return instance, nil
}